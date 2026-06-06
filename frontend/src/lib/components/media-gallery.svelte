<script lang="ts">
	import type { ClassValue } from "svelte/elements";

	import MediaGalleryCard from "$lib/components/media-gallery-card.svelte";
	import { mediaStore } from "$lib/stores/media-store.svelte";
	import { cn } from "$lib/utils";

	const { class: klass }: { class?: ClassValue } = $props();
</script>

<div
	class={cn(
		"grid grid-cols-2 gap-x-2 gap-y-4 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 2xl:grid-cols-6",
		klass,
	)}
>
	{#if mediaStore.isPreviewLoading}
		<div class="aspect-3/4 animate-pulse rounded bg-gray-500"></div>
		<div class="aspect-3/4 animate-pulse rounded bg-gray-500"></div>
		<div class="hidden aspect-3/4 animate-pulse rounded bg-gray-500 md:block"></div>
		<div class="hidden aspect-3/4 animate-pulse rounded bg-gray-500 lg:block"></div>
		<div class="hidden aspect-3/4 animate-pulse rounded bg-gray-500 xl:block"></div>
		<div class="hidden aspect-3/4 animate-pulse rounded bg-gray-500 2xl:block"></div>
	{:else}
		{#each mediaStore.mediaItems as media (media.pageUrl)}
			<MediaGalleryCard {media} />
		{/each}
	{/if}
</div>
