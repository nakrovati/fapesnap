import { AppService } from "$bindings/index";
import { CheckInterval, type DownloadDir } from "$bindings/internal/config";
import { parseWailsError } from "$lib/utils";
import { toast } from "svelte-sonner";

interface ConfigStore {
	autoCheck: boolean;
	checkIntervalValue: CheckInterval;
	checkOnStartup: boolean;
	downloadDirectory: DownloadDir | null;
	includePrereleases: boolean;
}

export const INTERVAL_OPTIONS = [
	{ value: CheckInterval.IntervalDay, label: "Day" },
	{ value: CheckInterval.Interval3Days, label: "3 days" },
	{ value: CheckInterval.IntervalWeek, label: "Week" },
	{ value: CheckInterval.IntervalMonth, label: "Month" },
] as const;

export const configStore = $state<ConfigStore>({
	checkIntervalValue: CheckInterval.IntervalWeek,
	checkOnStartup: false,
	downloadDirectory: null,
	includePrereleases: false,
	autoCheck: true,
});

export async function handleSelectDownloadDirectory() {
	try {
		configStore.downloadDirectory = await AppService.SelectDownloadDir();
	} catch (error) {
		toast.error("Error", { description: parseWailsError(error).message });
	}
}

export async function handleUnsetDownloadDir() {
	try {
		await AppService.UnsetDownloadDir();
		configStore.downloadDirectory = null;
	} catch (error) {
		toast.error("Error", { description: parseWailsError(error).message });
	}
}
