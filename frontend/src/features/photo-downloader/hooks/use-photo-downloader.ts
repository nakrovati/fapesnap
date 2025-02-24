import { DownloadPhotos, GetPhotos } from "$wails/go/main/App";
import * as wails from "$wails/runtime";
import { createEffect, createSignal, onCleanup } from "solid-js";
import { showToast } from "~/components/ui/toast";
import { photoStore, setPhotoStore } from "../stores/photo-store";

export function usePhotosDownloader() {
	const [collection, setCollection] = createSignal("");
	const [photos, setPhotos] = createSignal<string[]>(photoStore.photos);
	const [maxParallelDownloads, setMaxParallelDownloads] = createSignal(3);
	const [loading, setLoading] = createSignal(false);
	const [downloading, setDownloading] = createSignal(false);

	createEffect(() => {
		wails.EventsOn("download-start", () => {
			showToast({ title: "Download started" });
		});
		wails.EventsOn("download-complete", (description: string) => {
			showToast({
				title: "Download complete",
				description,
				variant: "success",
			});
		});
	});

	onCleanup(() => {
		wails.EventsOff("download-start");
		wails.EventsOff("download-complete");
	});

	function downloadPhotos(provider: string) {
		setDownloading(true);

		DownloadPhotos(collection(), provider, maxParallelDownloads())
			.catch((error) => {
				showToast({
					title: "Error",
					description: error,
					variant: "error",
				});
			})
			.finally(() => {
				setDownloading(false);
			});
	}

	function previewPhotos(provider: string) {
		setLoading(true);

		GetPhotos(collection(), provider)
			.then((result) => {
				setPhotos(result);

				setPhotoStore("photos", result);
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
		maxParallelDownloads,
		setMaxParallelDownloads,
		downloading,
		loading,
		downloadPhotos,
		previewPhotos,
	};
}
