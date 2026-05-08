import { toast } from "svelte-sonner";
import { DownloadPhoto, DownloadPhotos, GetPhotos } from "$lib/wailsjs/go/main/App";
import { providers } from "$lib/shared/constants";
import { providers as Providers } from "$lib/wailsjs/go/models";

type Photo = Providers.Photo & { providerName: string; collectionInput: string };

interface PhotoStore {
	providerName: string;
	collectionInput: string;
	photos: Photo[];
	maxParallelDownloads: string;
	loading: boolean;
	downloading: boolean;
}

export const photoStore = $state<PhotoStore>({
	providerName: providers[0]!.value,
	collectionInput: "",
	photos: [],
	maxParallelDownloads: "1",
	loading: false,
	downloading: false,
});

export function previewPhotos() {
	const { providerName, collectionInput } = photoStore;

	photoStore.loading = true;

	GetPhotos(collectionInput, providerName)
		.then((result) => {
			photoStore.photos = result.map((photo) => ({
				...photo,
				providerName,
				collectionInput,
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
	const { providerName, collectionInput, maxParallelDownloads } = photoStore;

	photoStore.downloading = true;

	DownloadPhotos(collectionInput, providerName, Number(maxParallelDownloads))
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
	const { providerName, collectionInput } = photoStore;

	DownloadPhoto(src, collectionInput, providerName)
		.then(() => {
			toast.success("Downloaded");
		})
		.catch((error) => {
			toast.error("Error", {
				description: error,
			});
		});
}
