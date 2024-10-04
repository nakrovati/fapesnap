import { DownloadPhotos, GetPhotos } from "$wails/go/main/App";
import { createSignal } from "solid-js";
import { showToast } from "~/components/ui/toast";
import { photoStore, setPhotoStore } from "../stores/photo-store";

export function usePhotosDownloader() {
  const [collection, setCollection] = createSignal("");
  const [photos, setPhotos] = createSignal<string[]>(photoStore.photos);
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

    GetPhotos(collection(), provider)
      .then((result) => {
        setPhotos(result);

        setPhotoStore("photos", result);

        console.log(result);
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
