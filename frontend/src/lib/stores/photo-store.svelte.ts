import { AppService } from "$bindings/index";
import { Photo } from "$bindings/internal/providers/models";
import { providers } from "$lib/shared/constants";
import { toast } from "svelte-sonner";

type ExtendedPhoto = Photo & { collectionInput: string; providerName: string };

interface PhotoStore {
	collectionInput: string;
	downloading: boolean;
	loading: boolean;
	maxParallelDownloads: string;
	photos: ExtendedPhoto[];
	providerName: string;
}

export const photoStore = $state<PhotoStore>({
	providerName: providers[0]!.value,
	collectionInput: "",
	photos: [],
	maxParallelDownloads: "1",
	loading: false,
	downloading: false,
});

export function downloadPhoto(src: string) {
	const { providerName, collectionInput } = photoStore;

	AppService.DownloadPhoto(src, collectionInput, providerName)
		.then(() => {
			toast.success("Downloaded");
		})
		.catch((error) => {
			toast.error("Error", {
				description: error,
			});
		});
}

export function downloadPhotos() {
	const { providerName, collectionInput, maxParallelDownloads } = photoStore;

	photoStore.downloading = true;

	AppService.DownloadPhotos(collectionInput, providerName, Number(maxParallelDownloads))
		.catch((error) => {
			toast.error("Error", {
				description: error,
			});
		})
		.finally(() => {
			photoStore.downloading = false;
		});
}

export function previewPhotos() {
	const { providerName, collectionInput } = photoStore;

	photoStore.loading = true;

	AppService.GetPhotos(collectionInput, providerName)
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
