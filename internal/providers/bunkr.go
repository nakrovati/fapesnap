package providers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

const bunkrUserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/148 Safari/537.36"

var (
	bunkrCDNPattern         = regexp.MustCompile(`var\s+jsCDN\s*=\s*"([^"]+)"`)
	bunkrSignPattern        = regexp.MustCompile(`var\s+signUrl\s*=\s*"([^"]+)"`)
	bunkrSignServicePattern = regexp.MustCompile(`SIGN_SERVICE_URL\s*=\s*["']([^"']+)["']`)
)

type bunkrSignResponse struct {
	Token string `json:"token"`
	Ex    int64  `json:"ex"`
}

type bunkrFileMetadata struct {
	MediaFiles string `json:"mediafiles"`
	Path       string `json:"path"`
	Original   string `json:"original"`
}

type BunkrProvider struct {
	ProviderName string
	BaseURL      string
}

func NewBunkrProvider() *BunkrProvider {
	return &BunkrProvider{
		ProviderName: "bunkr",
		BaseURL:      "https://bunkr.cr",
	}
}

func newBunkrCollector() *colly.Collector {
	c := colly.NewCollector()
	c.SetRequestTimeout(60 * time.Second)

	return c
}

func (p *BunkrProvider) FetchMediaItems(collectionSlug string) ([]Media, error) {
	mediaItems := make([]Media, 0)

	err := p.fetchAlbumPage(collectionSlug, &mediaItems)
	if err != nil {
		if errors.Is(err, ErrVisitNotFound) {
			return nil, fmt.Errorf("profile not found: %s", collectionSlug)
		}

		return []Media{}, err
	}

	if len(mediaItems) == 0 {
		return []Media{}, ErrNoMediaFound
	}

	return mediaItems, nil
}

func (p *BunkrProvider) GetCollectionFromURL(inputURL string) (string, error) {
	_, err := url.Parse(inputURL)
	if err != nil {
		return "", err
	}

	if !strings.Contains(inputURL, p.BaseURL) && !strings.Contains(inputURL, "bunkr") {
		return "", fmt.Errorf("%w: %s", ErrInvalidDomain, inputURL)
	}

	inputURL = strings.TrimSuffix(inputURL, "/")
	parts := strings.Split(inputURL, "/")

	if len(parts) < 5 || parts[len(parts)-1] == "" {
		return "", errors.New("can't get collection from url")
	}

	return parts[len(parts)-1], nil
}

func (p *BunkrProvider) GetMediaURL(pageURL string, collectionSlug string) (string, error) {
	c := newBunkrCollector()

	var (
		directURL       string
		cdnURL          string
		signURL         string
		fileDownloadURL string
		unavailableText string
		visitErr        error
	)

	// Some pages already contain a signed media URL.
	c.OnHTML("video source[src]", func(e *colly.HTMLElement) {
		directURL = e.Attr("src")
	})

	c.OnHTML("img#img-main[src]", func(e *colly.HTMLElement) {
		directURL = e.Attr("src")
	})

	c.OnHTML("main.cont img.w-full.h-full.absolute[src]", func(e *colly.HTMLElement) {
		directURL = e.Attr("src")
	})

	// Plain files use a separate download landing page.
	c.OnHTML("main.cont a[href*='/file/']", func(e *colly.HTMLElement) {
		fileDownloadURL = e.Attr("href")
	})

	// Bunkr disables downloads for files whose origin server is down.
	c.OnHTML("main.cont button[disabled]", func(e *colly.HTMLElement) {
		unavailableText = strings.TrimSpace(e.Text)
	})

	// Other pages only expose the raw CDN URL and a signing endpoint in JS.
	c.OnResponse(func(r *colly.Response) {
		cdnURL = pickBunkrScriptValue(bunkrCDNPattern, r.Body)
		signURL = pickBunkrScriptValue(bunkrSignPattern, r.Body)
	})

	c.OnError(func(c *colly.Response, err error) {
		visitErr = normalizeCollyError(c, err)
	})

	err := c.Visit(pageURL)

	if visitErr != nil {
		return "", visitErr
	}

	if err != nil {
		return "", err
	}

	if directURL != "" {
		return directURL, nil
	}

	if cdnURL != "" {
		return p.signMediaURL(pageURL, cdnURL, signURL)
	}

	if fileDownloadURL != "" {
		return p.getFileURL(pageURL, fileDownloadURL)
	}

	if unavailableText != "" {
		return "", fmt.Errorf("media unavailable: %s (%s)", pageURL, unavailableText)
	}

	return "", fmt.Errorf("media URL not found on page: %s", pageURL)
}

func (p *BunkrProvider) getFileURL(pageURL string, downloadURL string) (string, error) {
	baseURL, err := url.Parse(pageURL)
	if err != nil {
		return "", fmt.Errorf("parse page URL: %w", err)
	}

	landingURL, err := resolveBunkrURL(baseURL, downloadURL)
	if err != nil {
		return "", fmt.Errorf("parse download URL: %w", err)
	}

	var (
		fileID      string
		signService string
		visitErr    error
	)

	c := newBunkrCollector()

	c.OnHTML("#download-btn[data-id]", func(e *colly.HTMLElement) {
		fileID = e.Attr("data-id")
	})

	c.OnResponse(func(r *colly.Response) {
		signService = pickBunkrScriptValue(bunkrSignServicePattern, r.Body)
	})

	c.OnError(func(c *colly.Response, err error) {
		visitErr = normalizeCollyError(c, err)
	})

	err = c.Visit(landingURL.String())

	if visitErr != nil {
		return "", visitErr
	}

	if err != nil {
		return "", err
	}

	if fileID == "" {
		return "", fmt.Errorf("file ID not found on page: %s", landingURL)
	}

	if signService == "" {
		return "", fmt.Errorf("sign service not found on page: %s", landingURL)
	}

	metadata, err := p.getFileMetadata(landingURL, fileID)
	if err != nil {
		return "", err
	}

	rawURL := strings.TrimSuffix(metadata.MediaFiles, "/") + metadata.Path

	fileURL, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("parse file URL: %w", err)
	}

	if metadata.Original != "" {
		query := fileURL.Query()
		query.Set("n", metadata.Original)
		fileURL.RawQuery = query.Encode()
	}

	return p.signMediaURL(pageURL, fileURL.String(), signService)
}

func (p *BunkrProvider) getFileMetadata(landingURL *url.URL, fileID string) (bunkrFileMetadata, error) {
	metadataURL := *landingURL
	metadataURL.Path = "/api/_001_v2"
	metadataURL.RawQuery = ""
	metadataURL.Fragment = ""

	body, err := json.Marshal(map[string]string{"id": fileID})
	if err != nil {
		return bunkrFileMetadata{}, fmt.Errorf("marshal file metadata request: %w", err)
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, metadataURL.String(), bytes.NewReader(body))
	if err != nil {
		return bunkrFileMetadata{}, fmt.Errorf("build file metadata request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", bunkrUserAgent)
	req.Header.Set("Origin", landingURL.Scheme+"://"+landingURL.Host)
	req.Header.Set("Referer", landingURL.String())

	client := &http.Client{
		Transport:     http.DefaultTransport,
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return bunkrFileMetadata{}, fmt.Errorf("request file metadata: %w", err)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return bunkrFileMetadata{}, fmt.Errorf("request file metadata: unexpected status code: %d", resp.StatusCode)
	}

	responseBody, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return bunkrFileMetadata{}, fmt.Errorf("read file metadata: %w", err)
	}

	var metadata bunkrFileMetadata

	err = json.Unmarshal(responseBody, &metadata)
	if err != nil {
		return bunkrFileMetadata{}, fmt.Errorf("parse file metadata: %w", err)
	}

	if metadata.MediaFiles == "" || metadata.Path == "" {
		return bunkrFileMetadata{}, errors.New("file metadata is missing media URL")
	}

	return metadata, nil
}

func (p *BunkrProvider) signMediaURL(pageURL string, cdnURL string, signURL string) (string, error) {
	baseURL, err := url.Parse(pageURL)
	if err != nil {
		return "", fmt.Errorf("parse page URL: %w", err)
	}

	mediaURL, err := resolveBunkrURL(baseURL, cdnURL)
	if err != nil {
		return "", fmt.Errorf("parse media URL: %w", err)
	}

	signEndpoint, err := resolveBunkrURL(baseURL, signURL)
	if err != nil {
		return "", fmt.Errorf("parse sign URL: %w", err)
	}

	query := signEndpoint.Query()
	query.Set("path", mediaURL.Path)
	signEndpoint.RawQuery = query.Encode()

	signResponse, err := requestBunkrSignature(signEndpoint.String(), pageURL)
	if err != nil {
		return "", err
	}

	query = mediaURL.Query()
	query.Set("token", signResponse.Token)
	query.Set("ex", strconv.FormatInt(signResponse.Ex, 10))
	mediaURL.RawQuery = query.Encode()

	return mediaURL.String(), nil
}

func resolveBunkrURL(baseURL *url.URL, rawURL string) (*url.URL, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	if parsedURL.IsAbs() {
		return parsedURL, nil
	}

	return baseURL.ResolveReference(parsedURL), nil
}

func requestBunkrSignature(signEndpoint string, pageURL string) (bunkrSignResponse, error) {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, signEndpoint, nil)
	if err != nil {
		return bunkrSignResponse{}, fmt.Errorf("build sign request: %w", err)
	}

	req.Header.Set("User-Agent", bunkrUserAgent)
	req.Header.Set("Referer", pageURL)

	client := &http.Client{
		Transport:     http.DefaultTransport,
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return bunkrSignResponse{}, fmt.Errorf("sign media URL: %w", err)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return bunkrSignResponse{}, fmt.Errorf("sign media URL: unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return bunkrSignResponse{}, fmt.Errorf("read sign response: %w", err)
	}

	var signResponse bunkrSignResponse

	err = json.Unmarshal(body, &signResponse)
	if err != nil {
		return bunkrSignResponse{}, fmt.Errorf("parse sign response: %w", err)
	}

	if signResponse.Token == "" {
		return bunkrSignResponse{}, errors.New("sign response is missing token")
	}

	return signResponse, nil
}

func pickBunkrScriptValue(pattern *regexp.Regexp, body []byte) string {
	matches := pattern.FindSubmatch(body)
	if len(matches) < 2 {
		return ""
	}

	return strings.ReplaceAll(string(matches[1]), `\/`, "/")
}

func (p *BunkrProvider) fetchAlbumPage(collectionSlug string, mediaItems *[]Media) error {
	albumURL, err := url.JoinPath(p.BaseURL, "a", collectionSlug)
	if err != nil {
		return err
	}

	c := colly.NewCollector()

	var visitErr error

	c.OnHTML(".theItem", func(e *colly.HTMLElement) {
		item := p.parseItem(e)

		*mediaItems = append(*mediaItems, item)
	})

	c.OnError(func(c *colly.Response, err error) {
		visitErr = normalizeCollyError(c, err)
	})

	err = c.Visit(albumURL)

	if visitErr != nil {
		return visitErr
	}

	if err != nil {
		return err
	}

	return nil
}

func (p *BunkrProvider) parseItem(e *colly.HTMLElement) Media {
	mediaType := p.checkMediaType(e)

	href := e.ChildAttr("a[aria-label='download']", "href")

	pageURL, err := url.JoinPath(p.BaseURL, href)
	if err != nil {
		return Media{}
	}

	var thumbnailURL string
	if mediaType == MediaTypeImage || mediaType == MediaTypeVideo {
		thumbnailURL = e.ChildAttr("img.grid-images_box-img", "src")
	}

	return Media{
		PageURL:      pageURL,
		Name:         e.ChildText(".theName"),
		Size:         e.ChildText(".theSize"),
		ThumbnailURL: thumbnailURL,
		Type:         mediaType,
	}
}

func (p *BunkrProvider) checkMediaType(e *colly.HTMLElement) MediaType {
	class := e.ChildAttr(`span[class*="type-"]`, "class")

	switch {
	case strings.Contains(class, "type-Image"):
		return MediaTypeImage
	case strings.Contains(class, "type-Video"):
		return MediaTypeVideo
	case strings.Contains(class, "type-File"):
		return MediaTypeFile
	default:
		return MediaTypeUnknown
	}
}
