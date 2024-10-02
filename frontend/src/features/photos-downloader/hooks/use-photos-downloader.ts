import { DownloadPhotos, GetPhotos } from "$wails/go/main/App";
import { createSignal } from "solid-js";
import { showToast } from "~/components/ui/toast";

export function usePhotosDownloader() {
  const [collection, setCollection] = createSignal("");
  const [photos, setPhotos] = createSignal<string[]>([]);
  const [loading, setLoading] = createSignal(false);

  function downloadPhotos(provider: string) {
    DownloadPhotos(collection(), provider)
      .then((result) => {
        showToast({
          title: "Download has started",
          description: `${result.length} photos downloaded`,
          variant: "success",
        });
      })
      .catch((error) => {
        showToast({
          title: "Error",
          description: error,
          variant: "error",
        });
      });
  }

  function previewPhotos(provider: string) {
    setLoading(true);

    // const pattern = providerUrlPatterns[provider as Providers];
    // const result = processInput(collection(), pattern);

    GetPhotos(collection(), provider)
      .then((result) => {
        setPhotos(result.toReversed());
      })
      .catch((error) => {
        showToast({
          title: "Error",
          description: error,
          variant: "error",
        });
      })
      .finally(() => {
        setLoading(false);
      });
  }

  return {
    collection,
    setCollection,
    photos,
    loading,
    downloadPhotos,
    previewPhotos,
  };
}
