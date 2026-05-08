<script lang="ts">
	import { GetDownloadDir, UnsetDownloadDir, SelectDownloadDir } from "$lib/wailsjs/go/main/App";
	import ThemeSelect from "$lib/components/theme-select.svelte";
	import { Button } from "$lib/components/ui/button";
	import type { config } from "$lib/wailsjs/go/models";
	import X from "@lucide/svelte/icons/x";
	import Label from "$lib/components/ui/label/label.svelte";

	let selectedDownloadDir = $state<config.DownloadDir | null>();

	async function handleSelectDownloadDir() {
		const downloadDir = await SelectDownloadDir();
		selectedDownloadDir = downloadDir;
	}

	async function handleUnsetDownloadDir() {
		await UnsetDownloadDir();
		selectedDownloadDir = null;
	}

	$effect(() => {
		GetDownloadDir().then((dir) => {
			selectedDownloadDir = dir;
			console.log(dir);
		});
	});
</script>

<div class="flex flex-col gap-4">
	<div class="flex flex-col gap-1">
		<Label>Theme</Label>
		<ThemeSelect />
	</div>

	<div class="flex flex-col gap-1">
		<Label>Download directory</Label>
		<div class="flex items-center gap-2">
			<Button onclick={handleSelectDownloadDir}>Change directory</Button>

			{#if !!selectedDownloadDir && !selectedDownloadDir.isDefault}
				<Button
					onclick={handleUnsetDownloadDir}
					variant="outline"
					size="icon"
					aria-label="Unset directory"
				>
					<X />
				</Button>

				<span class="font-mono text-sm">{selectedDownloadDir.path}</span>
			{/if}
		</div>
	</div>
</div>
