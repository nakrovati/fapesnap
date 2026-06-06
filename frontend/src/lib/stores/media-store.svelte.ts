import { AppService } from "$bindings/index";
import { Media } from "$bindings/internal/providers/models";
import { providers } from "$lib/shared/constants";
import { parseWailsError } from "$lib/utils";
import { toast } from "svelte-sonner";

export type ExtendedMedia = Media & {
	collectionInput: string;
	downloadStatus?: DownloadStatus;
	providerName: string;
};

type DownloadStatus = "completed" | "downloading" | "failed" | "pending";

interface MediaStore {
	collectionInput: string;
	isDownloading: boolean;
	isPreviewLoading: boolean;
	maxParallelDownloads: string;
	mediaItems: ExtendedMedia[];
	providerName: string;
}

export const mediaStore = $state<MediaStore>({
	providerName: providers[0]!.value,
	collectionInput: "",
	mediaItems: [],
	maxParallelDownloads: "1",
	isPreviewLoading: false,
	isDownloading: false,
});

export async function downloadMedia(pageUrl: string) {
	const { providerName, collectionInput } = mediaStore;

	try {
		await AppService.DownloadMedia(pageUrl, collectionInput, providerName);

		toast.success("Downloaded");
	} catch (error) {
		toast.error("Error", {
			description: parseWailsError(error).message,
		});
	}
}

export async function downloadMediaItems() {
	const { providerName, collectionInput, maxParallelDownloads } = mediaStore;

	mediaStore.isDownloading = true;

	try {
		await AppService.DownloadMediaItems(
			collectionInput,
			providerName,
			Number(maxParallelDownloads),
		);
	} catch (error) {
		toast.error("Error", {
			description: parseWailsError(error).message,
		});
	} finally {
		mediaStore.isDownloading = false;
	}
}

export async function previewMediaItems() {
	const { providerName, collectionInput } = mediaStore;

	mediaStore.isPreviewLoading = true;

	try {
		const mediaItems = await AppService.GetMediaItems(collectionInput, providerName);

		mediaStore.mediaItems = mediaItems.map((media) => ({
			...media,
			providerName,
			collectionInput,
			downloadStatus: undefined,
		}));
	} catch (error) {
		toast.error("Error", {
			description: parseWailsError(error).message,
		});
	} finally {
		mediaStore.isPreviewLoading = false;
	}
}
