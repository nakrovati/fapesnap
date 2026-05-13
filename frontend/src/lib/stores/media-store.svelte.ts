import { AppService } from "$bindings/index";
import { Media } from "$bindings/internal/providers/models";
import { providers } from "$lib/shared/constants";
import { toast } from "svelte-sonner";

type ExtendedMedia = Media & { collectionInput: string; providerName: string };

interface MediaStore {
	collectionInput: string;
	downloading: boolean;
	loading: boolean;
	maxParallelDownloads: string;
	mediaItems: ExtendedMedia[];
	providerName: string;
}

export const mediaStore = $state<MediaStore>({
	providerName: providers[0]!.value,
	collectionInput: "",
	mediaItems: [],
	maxParallelDownloads: "1",
	loading: false,
	downloading: false,
});

export function downloadMedia(src: string) {
	const { providerName, collectionInput } = mediaStore;

	AppService.DownloadMedia(src, collectionInput, providerName)
		.then(() => {
			toast.success("Downloaded");
		})
		.catch((error) => {
			toast.error("Error", {
				description: error,
			});
		});
}

export function downloadMediaItems() {
	const { providerName, collectionInput, maxParallelDownloads } = mediaStore;

	mediaStore.downloading = true;

	AppService.DownloadMediaItems(collectionInput, providerName, Number(maxParallelDownloads))
		.catch((error) => {
			toast.error("Error", {
				description: error,
			});
		})
		.finally(() => {
			mediaStore.downloading = false;
		});
}

export function previewMediaItems() {
	const { providerName, collectionInput } = mediaStore;

	mediaStore.loading = true;

	AppService.GetMediaItems(collectionInput, providerName)
		.then((result) => {
			mediaStore.mediaItems = result.map((media) => ({
				...media,
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
			mediaStore.loading = false;
		});
}
