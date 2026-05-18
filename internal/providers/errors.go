package providers

import (
	"errors"
	"net/http"

	"github.com/gocolly/colly/v2"
)

var (
	ErrInvalidDomain = errors.New("invalid domain")
	ErrNoMediaFound  = errors.New("no media found")
	ErrVisitNotFound = errors.New("not found")
	ErrForbidden     = errors.New("forbidden")
)

func normalizeCollyError(r *colly.Response, err error) error {
	if r == nil {
		return err
	}

	switch r.StatusCode {
	case http.StatusNotFound:
		return ErrVisitNotFound
	case http.StatusForbidden:
		return ErrForbidden
	default:
		return err
	}
}
