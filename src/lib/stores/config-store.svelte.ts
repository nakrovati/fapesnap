import { toast } from "svelte-sonner";

import { parseError } from "#lib/utils.ts";
import { AppService, type CheckInterval, type DownloadDir } from "#mygo";

interface ConfigStore {
	autoCheck: boolean;
	checkIntervalValue: CheckInterval;
	downloadDirectory: DownloadDir | null;
	includePrereleases: boolean;
}

export const INTERVAL_OPTIONS = [
	{ value: 1440, label: "Day" },
	{ value: 4320, label: "3 days" },
	{ value: 10080, label: "Week" },
	{ value: 43200, label: "Month" },
] as const satisfies readonly { label: string; value: CheckInterval }[];

export const configStore = $state<ConfigStore>({
	checkIntervalValue: INTERVAL_OPTIONS[2].value,
	downloadDirectory: null,
	includePrereleases: false,
	autoCheck: true,
});

export async function handleSelectDownloadDirectory() {
	try {
		configStore.downloadDirectory = await AppService.selectDownloadDir();
	} catch (error) {
		toast.error("Error", { description: parseError(error).message });
	}
}

export async function handleUnsetDownloadDir() {
	try {
		await AppService.unsetDownloadDir();
		configStore.downloadDirectory = null;
	} catch (error) {
		toast.error("Error", { description: parseError(error).message });
	}
}
