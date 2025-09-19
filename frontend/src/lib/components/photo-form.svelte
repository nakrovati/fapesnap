<script lang="ts">
	import { Input } from "$lib/components/ui/input";
	import * as Select from "$lib/components/ui/select";
	import { providers } from "$lib/shared/constants";
	import { StopTask } from "$lib/wailsjs/go/main/App";
	import { photoStore, previewPhotos, downloadPhotos } from "$lib/stores/photo-store.svelte";
	import { Button } from "$lib/components/ui/button";
	import * as Tooltip from "$lib/components/ui/tooltip";

	let selectedProvider = $derived(providers.find((p) => p.value === photoStore.provider)!);
	let collectionTextFieldPlaceholder = $derived(
		selectedProvider.type === "id"
			? "Enter the album ID or URL"
			: "Enter the user's name or profile URL",
	);
</script>

<div class="mt-2">
	<div class="flex gap-2">
		<Input
			type="text"
			placeholder={collectionTextFieldPlaceholder}
			class="grow-1"
			autocorrect="off"
			bind:value={photoStore.collection}
		/>

		<Tooltip.Provider>
			<Tooltip.Root>
				<Tooltip.Trigger>
					<Select.Root bind:value={photoStore.maxParallelDownloads} type="single">
						<Select.Trigger>
							{photoStore.maxParallelDownloads}
						</Select.Trigger>
						<Select.Content>
							{#each [1, 2, 3, 4, 5] as n}
								<Select.Item value={n.toString()}>{n}</Select.Item>
							{/each}
						</Select.Content>
					</Select.Root>
				</Tooltip.Trigger>
				<Tooltip.Content>
					<p>Number of photos uploaded simultaneously</p>
				</Tooltip.Content>
			</Tooltip.Root>
		</Tooltip.Provider>

		<Select.Root bind:value={photoStore.provider} type="single">
			<Select.Trigger class="w-[180px]">
				{selectedProvider.label}
			</Select.Trigger>
			<Select.Content>
				{#each providers as provider (provider.value)}
					<Select.Item value={provider.value}>
						{provider.label}
					</Select.Item>
				{/each}
			</Select.Content>
		</Select.Root>
	</div>

	<div class="mt-4 flex justify-center gap-2">
		{#if !photoStore.loading}
			{#if photoStore.downloading}
				<Button onclick={StopTask} variant="destructive">Cancel</Button>
			{:else}
				<Button onclick={downloadPhotos}>Download</Button>
			{/if}
		{/if}
		{#if !photoStore.downloading}
			<Button onclick={previewPhotos}>Preview</Button>
		{/if}
	</div>
</div>
