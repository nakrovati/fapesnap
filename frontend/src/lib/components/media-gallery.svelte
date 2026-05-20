<script lang="ts">
	import type { ClassValue } from "svelte/elements";

	import { dev } from "$app/environment";
	import { Media, MediaType } from "$bindings/internal/providers/models";
	import { downloadMedia, mediaStore } from "$lib/stores/media-store.svelte";
	import { cn } from "$lib/utils";
	import Download from "@lucide/svelte/icons/download";
	import FolderArchive from "@lucide/svelte/icons/folder-archive";
	import Info from "@lucide/svelte/icons/info";
	import { Browser } from "@wailsio/runtime";

	import Badge from "./ui/badge/badge.svelte";
	import Button, { buttonVariants } from "./ui/button/button.svelte";
	import * as Tooltip from "./ui/tooltip";

	const { class: klass }: { class?: ClassValue } = $props();

	async function handleDownloadMedia(media: Media) {
		if (
			media.type === MediaType.MediaTypeVideo ||
			media.type === MediaType.MediaTypeFile ||
			media.type === MediaType.MediaTypeUnknown
		) {
			await Browser.OpenURL(media.pageUrl);
		} else {
			downloadMedia(media.pageUrl);
		}
	}
</script>

<div
	class={cn(
		"grid grid-cols-2 gap-x-2 gap-y-4 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 2xl:grid-cols-6",
		klass,
	)}
>
	{#if mediaStore.loading}
		<div class="aspect-3/4 animate-pulse rounded bg-gray-500"></div>
		<div class="aspect-3/4 animate-pulse rounded bg-gray-500"></div>
		<div class="hidden aspect-3/4 animate-pulse rounded bg-gray-500 md:block"></div>
		<div class="hidden aspect-3/4 animate-pulse rounded bg-gray-500 lg:block"></div>
		<div class="hidden aspect-3/4 animate-pulse rounded bg-gray-500 xl:block"></div>
		<div class="hidden aspect-3/4 animate-pulse rounded bg-gray-500 2xl:block"></div>
	{:else}
		{#each mediaStore.mediaItems as media (media.pageUrl)}
			<div class="relative">
				<div class="absolute top-1 right-1 flex items-center gap-2">
					{#if media.type === MediaType.MediaTypeVideo || media.type === MediaType.MediaTypeFile || media.type === MediaType.MediaTypeUnknown}
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
												<dd class=" font-mono break-all text-zinc-800">{value}</dd>
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
						onclick={() => handleDownloadMedia(media)}
					>
						<Download />
					</Button>
				</div>
				<div class="flex w-full justify-center">
					{#if media.type === MediaType.MediaTypeImage || media.type === MediaType.MediaTypeVideo}
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
		{/each}
	{/if}
</div>
