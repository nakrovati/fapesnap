import { AppService, UpdateService } from "#mygo";

import type { PageLoad } from "./$types";

export const load: PageLoad = async () => {
	const [downloadDir, updates, updatesEnabled] = await Promise.all([
		AppService.getDownloadDir(),
		UpdateService.getUpdatesConfig(),
		UpdateService.updatesEnabled(),
	]);
	return {
		downloadDir,
		updates,
		updatesEnabled,
		meta: {
			title: "Settings",
		},
	};
};
