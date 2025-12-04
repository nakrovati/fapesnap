import { toast } from "svelte-sonner";
import { DownloadPhoto, DownloadPhotos, GetPhotos } from "$lib/wailsjs/go/main/App";
import { providers } from "$lib/shared/constants";
import { providers as Providers } from "$lib/wailsjs/go/models";

type Photo = Providers.Photo & { provider: string; collection: string };

interface PhotoStore {
	provider: string;
	collection: string;
	photos: Photo[];
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
	const { provider, collection } = photoStore;

	photoStore.loading = true;

	GetPhotos(collection, provider)
		.then((result) => {
			photoStore.photos = result.map((photo) => ({
				...photo,
				provider,
				collection,
			}));
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
	const { provider, collection, maxParallelDownloads } = photoStore;

	photoStore.downloading = true;

	DownloadPhotos(collection, provider, Number(maxParallelDownloads))
		.catch((error) => {
			toast.error("Error", {
				description: error,
			});
		})
		.finally(() => {
			photoStore.downloading = false;
		});
}

export function downloadPhoto(src: string) {
	const { provider, collection } = photoStore;

	DownloadPhoto(src, collection, provider)
		.then(() => {
			toast.success("Downloaded");
		})
		.catch((error) => {
			toast.error("Error", {
				description: error,
			});
		});
}
