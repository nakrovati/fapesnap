import { Button } from "~/components/ui/button";
import { TextField, TextFieldInput } from "~/components/ui/text-field";
import {
  PhotosPreview,
  usePhotosDownloader,
} from "~/features/photos-downloader";
import {
  ProvidersSelector,
  useProviderSelector,
} from "~/features/providers-selector";

import { StopTask } from "$wails/go/main/App";

function IndexPage() {
  const { provider, providers, setProvider } = useProviderSelector();
  const { setCollection, photos, loading, previewPhotos, downloadPhotos } =
    usePhotosDownloader();

  const collectionTextFieldPlaceholder = () =>
    provider().type === "id"
      ? "Enter the album ID or URL"
      : "Enter the user's name or profile URL";

  function handleDownloadPhotos() {
    downloadPhotos(provider().value);
  }

  function handlePreviewPhotos() {
    previewPhotos(provider().value);
  }

  function handleCancelTask() {
    StopTask().then();
  }

  return (
    <div class="flex flex-col grow h-full">
      <div class="flex gap-2">
        <TextField class="grow">
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
        <Button onClick={handleCancelTask} variant="secondary">
          Cancel
        </Button>
      </div>

      <div class="overflow-y-auto mt-4">
        <PhotosPreview photos={photos()} loading={loading()} />
      </div>
    </div>
  );
}

export default IndexPage;
