<script lang="ts">
	import type { ClassValue } from "svelte/elements";

	import { MediaType } from "$bindings/internal/providers/models";
	import { downloadMedia, mediaStore } from "$lib/stores/media-store.svelte";
	import { cn } from "$lib/utils";
	import Download from "@lucide/svelte/icons/download";
	import FolderArchive from "@lucide/svelte/icons/folder-archive";

	import Badge from "./ui/badge/badge.svelte";
	import Button from "./ui/button/button.svelte";

	const { class: klass }: { class?: ClassValue } = $props();
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
		{#each mediaStore.mediaItems as media (media.url)}
			<div class="relative">
				<div class="absolute top-1 right-1 flex items-center gap-2">
					{#if media.type === MediaType.MediaTypeVideo || media.type === MediaType.MediaTypeFile}
						<Badge>{media.type}</Badge>
					{/if}

					<Button
						aria-label="Download media"
						size="icon"
						class="size-8"
						onclick={() => downloadMedia(media.url)}
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
					<p class="text-center">{media.name}</p>
				{/if}
				{#if media.size}
					<p class="text-center text-sm text-zinc-500">{media.size}</p>
				{/if}
			</div>
		{/each}
	{/if}
</div>
