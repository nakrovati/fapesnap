import { Button } from "~/components/ui/button";
import {
  TextField,
  TextFieldInput,
  TextFieldLabel,
} from "~/components/ui/text-field";
import {
  PhotosPreview,
  usePhotosDownloader,
} from "~/features/photos-downloader";
import {
  ProvidersSelector,
  useProviderSelector,
} from "~/features/providers-selector";

function IndexPage() {
  const { provider, providers, setProvider } = useProviderSelector();
  const { setCollection, photos, loading, previewPhotos, downloadPhotos } =
    usePhotosDownloader();

  const collectionTextFieldPlaceholder = () =>
    provider().type === "id"
      ? "Enter the album ID or URL"
      : "Enter the user's name or profile URL";
  const collectionTextFieldLabel = () =>
    provider().type === "id" ? "Album ID" : "Username";

  function handleDownloadPhotos() {
    downloadPhotos(provider().value);
  }

  function handlePreviewPhotos() {
    previewPhotos(provider().value);
  }

  return (
    <>
      <div class="flex place-items-end gap-2">
        <TextField class="w-full">
          <TextFieldLabel for="collection">
            {collectionTextFieldLabel()}
          </TextFieldLabel>
          <TextFieldInput
            placeholder={collectionTextFieldPlaceholder()}
            type="text"
            autocomplete="off"
            id="collection"
            onInput={(e) => setCollection(e.currentTarget.value)}
          />
        </TextField>

        <ProvidersSelector
          providers={providers}
          provider={provider()}
          onChange={setProvider}
        />
      </div>

      <div class="w-full flex justify-center gap-2 mt-4">
        <Button onClick={handleDownloadPhotos}>Download</Button>
        <Button onClick={handlePreviewPhotos} variant="secondary">
          Preview
        </Button>
      </div>

      <PhotosPreview photos={photos()} loading={loading()} />
    </>
  );
}

export default IndexPage;
