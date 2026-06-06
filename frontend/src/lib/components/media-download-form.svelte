<script lang="ts">
	import { StopTask } from "$bindings/appservice";
	import { Button } from "$lib/components/ui/button";
	import * as ButtonGroup from "$lib/components/ui/button-group";
	import { Input } from "$lib/components/ui/input";
	import * as Select from "$lib/components/ui/select";
	import * as Tooltip from "$lib/components/ui/tooltip";
	import { providers } from "$lib/shared/constants";
	import {
		downloadMediaItems,
		mediaStore,
		previewMediaItems,
	} from "$lib/stores/media-store.svelte";
	import { Events } from "@wailsio/runtime";
	import { toast } from "svelte-sonner";

	let selectedProvider = $derived(providers.find((p) => p.value === mediaStore.providerName)!);
	let collectionTextFieldPlaceholder = $derived(
		selectedProvider.type === "id"
			? "Enter the album ID or URL"
			: "Enter the user's name or profile URL",
	);

	const LAST_SELECTED_PROVIDER_KEY = "last-selected-provider";

	$effect(() => {
		mediaStore.providerName =
			localStorage.getItem(LAST_SELECTED_PROVIDER_KEY) ?? providers[0]!.value;
		localStorage.setItem(LAST_SELECTED_PROVIDER_KEY, mediaStore.providerName);
	});

	function handleProviderChange() {
		localStorage.setItem(LAST_SELECTED_PROVIDER_KEY, mediaStore.providerName);
	}

	$effect(() => {
		Events.On("download:batch-started", () => {
			for (const item of mediaStore.mediaItems) {
				item.downloadStatus = "pending";
			}

			toast("Download started");
		});
		Events.On("download:batch-completed", (event) => {
			toast.success("Download complete", {
				description: `Downloading completed! ${event.data.downloaded} of ${event.data.total} items downloaded.`,
			});
		});

		Events.On("download:media-started", ({ data }) => {
			const mediaItem = mediaStore.mediaItems.find((m) => m.pageUrl === data.pageUrl);
			if (mediaItem) mediaItem.downloadStatus = "downloading";
		});
		Events.On("download:media-completed", ({ data }) => {
			const mediaItem = mediaStore.mediaItems.find((m) => m.pageUrl === data.pageUrl);
			if (mediaItem) mediaItem.downloadStatus = "completed";
		});
		Events.On("download:media-failed", ({ data }) => {
			const mediaItem = mediaStore.mediaItems.find((m) => m.pageUrl === data.pageUrl);
			if (mediaItem) mediaItem.downloadStatus = "failed";
		});

		return () => {
			Events.Off(
				"download:batch-started",
				"download:batch-completed",
				"download:media-started",
				"download:media-completed",
				"download:media-failed",
			);
		};
	});

	async function handleCancelCurrentTasks() {
		await StopTask();

		for (const item of mediaStore.mediaItems) {
			if (item.downloadStatus === "pending") {
				item.downloadStatus = undefined;
			}
		}
	}
</script>

<div class="mt-2">
	<div class="flex gap-2">
		<Input
			type="text"
			placeholder={collectionTextFieldPlaceholder}
			class="grow"
			autocorrect="off"
			bind:value={mediaStore.collectionInput}
		/>
		<ButtonGroup.Root>
			<Select.Root
				bind:value={mediaStore.providerName}
				onValueChange={handleProviderChange}
				type="single"
			>
				<Select.Trigger class="w-45">
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

			<Tooltip.Provider>
				<Tooltip.Root>
					<Tooltip.Trigger>
						<Select.Root bind:value={mediaStore.maxParallelDownloads} type="single">
							<Select.Trigger class="rounded-l-none">
								{mediaStore.maxParallelDownloads}
							</Select.Trigger>
							<Select.Content>
								{#each [1, 2, 3, 4, 5] as n (n)}
									<Select.Item value={n.toString()}>{n}</Select.Item>
								{/each}
							</Select.Content>
						</Select.Root>
					</Tooltip.Trigger>
					<Tooltip.Content>
						<p>Number of media downloaded simultaneously</p>
					</Tooltip.Content>
				</Tooltip.Root>
			</Tooltip.Provider>
		</ButtonGroup.Root>
	</div>

	<div class="mt-4 flex justify-center gap-2">
		{#if !mediaStore.isPreviewLoading}
			{#if mediaStore.isDownloading}
				<Button onclick={handleCancelCurrentTasks} variant="destructive">Cancel</Button>
			{:else}
				<Button onclick={downloadMediaItems}>Download all</Button>
			{/if}
		{/if}
		{#if !mediaStore.isDownloading}
			<Button onclick={previewMediaItems} variant="secondary">Preview</Button>
		{/if}
	</div>
</div>
