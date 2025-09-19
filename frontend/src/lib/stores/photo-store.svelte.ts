import { toast } from "svelte-sonner";
import { DownloadPhotos, GetPhotos } from "$lib/wailsjs/go/main/App";
import { providers } from "$lib/shared/constants";

interface PhotoStore {
	provider: string;
	collection: string;
	photos: string[];
	maxParallelDownloads: string;
	loading: boolean;
	downloading: boolean;
}

export const photoStore = $state<PhotoStore>({
	provider: providers[0]!.value,
	collection: "",
	photos: [],
	maxParallelDownloads: "3",
	loading: false,
	downloading: false,
});

export function previewPhotos() {
	photoStore.loading = true;

	GetPhotos(photoStore.collection, photoStore.provider)
		.then((result) => {
			photoStore.photos = result;
			console.log("photos", photoStore.photos);
		})
		.catch((error) => {
			toast.error("Error", {
				description: error,
			});
		})
		.finally(() => {
			photoStore.loading = false;
		});
}

export function downloadPhotos() {
	photoStore.downloading = true;

	DownloadPhotos(
		photoStore.collection,
		photoStore.provider,
		Number(photoStore.maxParallelDownloads),
	)
		.catch((error) => {
			toast.error("Error", {
				description: error,
			});
		})
		.finally(() => {
			photoStore.downloading = false;
		});
}
