<script lang="ts">
	import { photoStore } from "$lib/stores/photo-store.svelte";
	import type { ClassValue } from "svelte/elements";

	const { class: klass }: { class?: ClassValue } = $props();
</script>

<div
	class={[
		"grid grid-cols-2 gap-x-2 gap-y-4 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 2xl:grid-cols-6",
		klass,
	]}
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
			<div>
				<img src={photo} alt="" class="min-h-48 rounded object-contain object-top" loading="lazy" />
			</div>
		{/each}
	{/if}
</div>
