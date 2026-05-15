<script lang="ts">
	import type { ClassValue } from "svelte/elements";

	import { MediaType } from "$bindings/internal/providers/models";
	import { downloadMedia, mediaStore } from "$lib/stores/media-store.svelte";
	import { cn } from "$lib/utils";
	import Download from "@lucide/svelte/icons/download";

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
					{#if media.type === MediaType.MediaTypeVideo}
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
				<img
					src={media.thumbnailUrl ?? media.url}
					alt=""
					class="min-h-48 rounded object-contain object-top"
					loading="lazy"
				/>
			</div>
		{/each}
	{/if}
</div>
