import { AppService, UpdateService } from "$bindings/index";

import type { PageLoad } from "./$types";

export const load: PageLoad = async () => {
	const [downloadDir, updates] = await Promise.all([
		AppService.GetDownloadDir(),
		UpdateService.GetUpdatesConfig(),
	]);
	return {
		downloadDir,
		updates,
		meta: {
			title: "Settings",
		},
	};
};
