<script lang="ts">
	import { downloadPhoto, photoStore } from "$lib/stores/photo-store.svelte";
	import { cn } from "$lib/utils";
	import type { ClassValue } from "svelte/elements";
	import Button from "./ui/button/button.svelte";
	import { Download } from "@lucide/svelte";

	const { class: klass }: { class?: ClassValue } = $props();
</script>

<div
	class={cn(
		"grid grid-cols-2 gap-x-2 gap-y-4 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 2xl:grid-cols-6",
		klass,
	)}
>
	{#if photoStore.loading}
		<div class="aspect-[3/4] animate-pulse rounded bg-gray-500"></div>
		<div class="aspect-[3/4] animate-pulse rounded bg-gray-500"></div>
		<div class="hidden aspect-[3/4] animate-pulse rounded bg-gray-500 md:block"></div>
		<div class="hidden aspect-[3/4] animate-pulse rounded bg-gray-500 lg:block"></div>
		<div class="hidden aspect-[3/4] animate-pulse rounded bg-gray-500 xl:block"></div>
		<div class="hidden aspect-[3/4] animate-pulse rounded bg-gray-500 2xl:block"></div>
	{:else}
		{#each photoStore.photos as photo}
			<div class="relative">
				<Button
					aria-label="Download photo"
					size="icon"
					class="absolute right-1 top-1 size-8"
					onclick={() => downloadPhoto(photo.url)}
				>
					<Download />
				</Button>
				<img
					src={photo.thumbnailUrl ?? photo.url}
					alt=""
					class="min-h-48 rounded object-contain object-top"
					loading="lazy"
				/>
			</div>
		{/each}
	{/if}
</div>
