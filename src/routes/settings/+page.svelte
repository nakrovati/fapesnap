<script lang="ts">
	import FolderOpen from "@lucide/svelte/icons/folder-open";
	import RefreshCw from "@lucide/svelte/icons/refresh-cw";
	import X from "@lucide/svelte/icons/x";
	import { toast } from "svelte-sonner";

	import ThemeSelect from "#lib/components/theme-select.svelte";
	import { Button } from "#lib/components/ui/button/index.ts";
	import { Label } from "#lib/components/ui/label/index.ts";
	import * as Select from "#lib/components/ui/select/index.ts";
	import { Separator } from "#lib/components/ui/separator/index.ts";
	import { Switch } from "#lib/components/ui/switch/index.ts";
	import {
		configStore,
		handleSelectDownloadDirectory,
		handleUnsetDownloadDir,
		INTERVAL_OPTIONS,
	} from "#lib/stores/config-store.svelte.ts";
	import { parseError } from "#lib/utils.ts";
	import { UpdateService } from "#mygo";

	import type { PageProps } from "./$types";

	let { data }: PageProps = $props();

	let statusMessage = $state("");

	let selectedCheckInterval = $derived(String(configStore.checkIntervalValue));
	let selectedIncludePrerelease = $derived(String(configStore.includePrereleases));

	let intervalLabel = $derived(
		INTERVAL_OPTIONS.find((o) => o.value === configStore.checkIntervalValue)?.label ??
			"Select interval",
	);

	let channelLabel = $derived(configStore.includePrereleases ? "Prerelease" : "Stable");

	$effect(() => {
		loadSettings();
	});

	async function handleAutoCheckChange(value: boolean) {
		statusMessage = "";

		try {
			const updates = await UpdateService.getUpdatesConfig();
			await UpdateService.setUpdatesConfig({
				...updates,
				autoCheck: value,
			});
			configStore.autoCheck = value;
			statusMessage = value ? "Auto-check enabled." : "Auto-check disabled.";
		} catch (error) {
			toast.error("Error", {
				description: parseError(error).message,
			});
			statusMessage = "Could not save update settings.";
		}
	}

	async function handleChannelChange(value: string) {
		const isPrerelease = Boolean(value);
		statusMessage = "Updating update channel...";

		try {
			const updates = await UpdateService.getUpdatesConfig();
			await UpdateService.setUpdatesConfig({
				...updates,
				includePrereleases: isPrerelease,
			});
			configStore.includePrereleases = isPrerelease;
			statusMessage = "Update channel changed.";
		} catch (error) {
			toast.error("Error", {
				description: parseError(error).message,
			});
			statusMessage = "Could not change update channel.";
		}
	}

	async function handleIntervalChange(value: string) {
		const checkInterval = INTERVAL_OPTIONS.find(
			(option) => option.value.toString() === value,
		)?.value;
		if (!checkInterval) {
			toast.error("Error", { description: "Invalid update check interval." });
			return;
		}

		statusMessage = "";

		try {
			const updates = await UpdateService.getUpdatesConfig();
			await UpdateService.setUpdatesConfig({
				...updates,
				checkIntervalMinutes: checkInterval,
			});
			configStore.checkIntervalValue = checkInterval;
			statusMessage = "Check interval saved. Restart the app to apply it.";
		} catch (error) {
			toast.error("Error", {
				description: parseError(error).message,
			});
			statusMessage = "Could not save update interval.";
		}
	}

	async function handleManualCheck() {
		statusMessage = "Opening update checker...";

		try {
			await UpdateService.check();
			statusMessage = "Update checker opened.";
		} catch (error) {
			toast.error("Error", {
				description: parseError(error).message,
			});
			statusMessage = "Update check failed.";
		}
	}

	function loadSettings() {
		configStore.downloadDirectory = data.downloadDir;
		configStore.autoCheck = data.updates.autoCheck;
		configStore.checkIntervalValue = data.updates.checkIntervalMinutes;
		configStore.includePrereleases = data.updates.includePrereleases;
	}
</script>

<div class="mx-auto flex max-w-3xl flex-col gap-8 pb-8">
	<!-- Appearance -->
	<section class="flex flex-col gap-3">
		<h3 class="px-1 text-sm font-medium text-muted-foreground">Appearance</h3>

		<div class="rounded-lg border bg-card p-4">
			<div class="flex items-center justify-between gap-4">
				<Label>Theme</Label>
				<ThemeSelect />
			</div>
		</div>
	</section>

	<!-- Downloads -->
	<section class="flex flex-col gap-3">
		<h3 class="px-1 text-sm font-medium text-muted-foreground">Downloads</h3>

		<div class="rounded-lg border bg-card p-4">
			<div class="flex flex-col gap-2">
				<Label>Download directory</Label>
				<div class="flex items-center gap-2">
					<div
						class="flex h-9 min-w-0 flex-1 items-center gap-2 rounded-md border bg-muted/40 px-3 text-sm"
					>
						<FolderOpen class="size-4 shrink-0 text-muted-foreground" />
						<span class="truncate font-mono text-xs">
							{#if !!configStore.downloadDirectory && !configStore.downloadDirectory.isDefault}
								{configStore.downloadDirectory.path}
							{:else}
								<span class="font-sans text-muted-foreground italic">Default directory</span>
							{/if}
						</span>
					</div>
					<Button variant="secondary" size="sm" class="h-9" onclick={handleSelectDownloadDirectory}>
						Change
					</Button>

					{#if !!configStore.downloadDirectory && !configStore.downloadDirectory.isDefault}
						<Button
							onclick={handleUnsetDownloadDir}
							variant="outline"
							size="icon"
							class="size-9 shrink-0"
							aria-label="Unset directory"
						>
							<X class="size-4" />
						</Button>
					{/if}
				</div>
			</div>
		</div>
	</section>

	<!-- Updates -->
	<section class="flex flex-col gap-3">
		<div class="flex items-center justify-between px-1">
			<h3 class="text-sm font-medium text-muted-foreground">Updates</h3>
			<div class="flex items-center gap-2">
				{#if statusMessage}
					<span class="animate-in text-xs text-muted-foreground fade-in">{statusMessage}</span>
				{/if}
				<Button
					onclick={handleManualCheck}
					variant="outline"
					size="sm"
					class="h-7 gap-1.5 text-xs"
					disabled={!data.updatesEnabled}
				>
					<RefreshCw class="size-3.5" />
					<span>Check now</span>
				</Button>
			</div>
		</div>
		{#if !data.updatesEnabled}
			<p class="px-1 text-xs text-muted-foreground">
				Self-updates are unavailable in development builds, builds without an update signing key, or
				installations that cannot be updated in place.
			</p>
		{/if}

		<div class="flex flex-col gap-4 rounded-lg border bg-card p-4">
			<div class="flex items-center justify-between gap-4">
				<Label for="auto-check">Enable automatic update checks</Label>
				<Switch
					id="auto-check"
					bind:checked={configStore.autoCheck}
					onCheckedChange={handleAutoCheckChange}
				/>
			</div>

			{#if configStore.autoCheck}
				<Separator />

				<div class="flex items-center justify-between gap-4">
					<div class="flex flex-col gap-1">
						<Label>Check every</Label>
						<span class="text-xs text-muted-foreground">Changes take effect after restart.</span>
					</div>
					<Select.Root
						value={selectedCheckInterval}
						type="single"
						onValueChange={handleIntervalChange}
					>
						<Select.Trigger class="h-9 w-40">
							{intervalLabel}
						</Select.Trigger>
						<Select.Content>
							{#each INTERVAL_OPTIONS as option (option.value)}
								<Select.Item value={option.value.toString()}>{option.label}</Select.Item>
							{/each}
						</Select.Content>
					</Select.Root>
				</div>
			{/if}

			<Separator />

			<div class="flex items-center justify-between gap-4">
				<div class="flex items-center gap-2">
					<Label>Update channel</Label>
				</div>
				<Select.Root
					value={selectedIncludePrerelease}
					type="single"
					onValueChange={handleChannelChange}
				>
					<Select.Trigger class="h-9 w-40">
						{channelLabel}
					</Select.Trigger>
					<Select.Content>
						<Select.Item value="false">Stable</Select.Item>
						<Select.Item value="true">Prerelease</Select.Item>
					</Select.Content>
				</Select.Root>
			</div>
		</div>
	</section>
</div>
