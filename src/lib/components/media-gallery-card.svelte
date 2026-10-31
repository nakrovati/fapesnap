<script lang="ts">
	import { dev } from "$app/env";
	import Check from "@lucide/svelte/icons/check";
	import Clock3 from "@lucide/svelte/icons/clock-3";
	import Download from "@lucide/svelte/icons/download";
	import FolderArchive from "@lucide/svelte/icons/folder-archive";
	import Info from "@lucide/svelte/icons/info";
	import LoaderCircle from "@lucide/svelte/icons/loader-circle";
	import X from "@lucide/svelte/icons/x";
	import { toast } from "svelte-sonner";

	import { Badge } from "#lib/components/ui/badge/index.ts";
	import { Button, buttonVariants } from "#lib/components/ui/button/index.ts";
	import * as Tooltip from "#lib/components/ui/tooltip/index.ts";
	import { downloadMedia, type ExtendedMedia } from "#lib/stores/media-store.svelte.ts";
	import { parseError } from "#lib/utils.ts";
	import { AppService } from "#mygo";

	let { media }: { media: ExtendedMedia } = $props();

	let status = $derived(media.downloadStatus);

	async function handleDownloadMedia() {
		if (
			media.providerName === "bunkr" &&
			(media.type === "video" || media.type === "file" || media.type === "unknown")
		) {
			try {
				await AppService.openExternal(media.pageUrl);
			} catch (error) {
				toast.error("Error", { description: parseError(error).message });
			}
		} else {
			await downloadMedia(media.pageUrl);
		}
	}
</script>

<div class="relative">
	<div class="absolute top-1 right-1 flex items-center gap-2">
		{#if media.type === "video" || media.type === "file" || media.type === "unknown"}
			<Badge>{media.type}</Badge>
		{/if}

		{#if dev}<Tooltip.Provider>
				<Tooltip.Root>
					<Tooltip.Trigger class={buttonVariants({ variant: "secondary", size: "icon" })}>
						<Info />
					</Tooltip.Trigger>
					<Tooltip.Content>
						<div class="flex flex-col gap-1">
							{#each Object.entries(media) as [key, value] (key)}
								<dl class="flex gap-1">
									<dt>{key}:</dt>
									<dd class="font-mono break-all">{value}</dd>
								</dl>
							{/each}
						</div>
					</Tooltip.Content>
				</Tooltip.Root>
			</Tooltip.Provider>
		{/if}

		<Button
			aria-label="Download media"
			size="icon"
			class="size-8"
			disabled={status === "downloading"}
			onclick={handleDownloadMedia}
		>
			{#if status === "completed"}
				<Check class="text-green-500" />
			{:else if status === "pending"}
				<Clock3 class="animate-pulse" />
			{:else if status === "failed"}
				<X class="text-red-500" />
			{:else if status === "downloading"}
				<LoaderCircle class="animate-spin text-blue-500" />
			{:else}
				<Download />
			{/if}
		</Button>
	</div>
	<div class="flex w-full justify-center">
		{#if media.type === "image" || media.type === "video"}
			<img
				src={media.thumbnailUrl ?? media.url}
				alt=""
				class="min-h-16 rounded object-contain object-top"
				loading="lazy"
			/>
		{:else}
			<div>
				<FolderArchive class="size-36 text-zinc-700" />
			</div>
		{/if}
	</div>
	{#if media.name}
		<p class="line-clamp-2 w-full text-center">{media.name}</p>
	{/if}
	{#if media.size}
		<p class="text-center text-sm text-zinc-500">{media.size}</p>
	{/if}
</div>
