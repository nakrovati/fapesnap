import { toast } from "svelte-sonner";
import { AppService } from "$bindings/index";
import { providers } from "$lib/shared/constants";
import { Photo } from "$bindings/internal/providers/models";

type ExtendedPhoto = Photo & { providerName: string; collectionInput: string };

interface PhotoStore {
	providerName: string;
	collectionInput: string;
	photos: ExtendedPhoto[];
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
