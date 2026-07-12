<script lang="ts">
	import { UpdateService } from "$bindings/index";
	import { CheckInterval } from "$bindings/internal/config";
	import ThemeSelect from "$lib/components/theme-select.svelte";
	import { Badge } from "$lib/components/ui/badge";
	import { Button } from "$lib/components/ui/button";
	import { Label } from "$lib/components/ui/label";
	import * as Select from "$lib/components/ui/select";
	import { Separator } from "$lib/components/ui/separator";
	import { Switch } from "$lib/components/ui/switch";
	import {
		configStore,
		handleSelectDownloadDirectory,
		handleUnsetDownloadDir,
	} from "$lib/stores/config-store.svelte";
	import { parseWailsError } from "$lib/utils";
	import FolderOpen from "@lucide/svelte/icons/folder-open";
	import RefreshCw from "@lucide/svelte/icons/refresh-cw";
	import X from "@lucide/svelte/icons/x";
	import { toast } from "svelte-sonner";

	import type { PageProps } from "./$types";

	let { data }: PageProps = $props();

	const INTERVAL_OPTIONS = [
		{ value: CheckInterval.IntervalDay, label: "Day" },
		{ value: CheckInterval.Interval3Days, label: "3 days" },
		{ value: CheckInterval.IntervalWeek, label: "Week" },
		{ value: CheckInterval.IntervalMonth, label: "Month" },
	] as const;

	let statusMessage = $state("");

	let selectedCheckInterval = $derived(String(configStore.checkIntervalValue));
	let selectedIncludePrerelease = $derived(String(configStore.includePrereleases));

	let intervalLabel = $derived(
		INTERVAL_OPTIONS.find((o) => o.value === +configStore.checkIntervalValue)?.label ??
			"Select interval",
	);

	let channelLabel = $derived(configStore.includePrereleases ? "Prerelease" : "Stable");

	$effect(() => {
		loadSettings();
	});

	function loadSettings() {
		configStore.downloadDirectory = data.downloadDir;
		configStore.autoCheck = data.updates.autoCheck;
		configStore.checkOnStartup = data.updates.checkOnStartup;
		configStore.checkIntervalValue = data.updates.checkIntervalMinutes;
		configStore.includePrereleases = data.updates.includePrereleases;
	}

	async function handleAutoCheckChange(value: boolean) {
		statusMessage = "";

		try {
			const updates = await UpdateService.GetUpdatesConfig();
			await UpdateService.SetUpdatesConfig({
				...updates,
				autoCheck: value,
			});
			configStore.autoCheck = value;
			statusMessage = value ? "Auto-check enabled." : "Auto-check disabled.";
		} catch (error) {
			toast.error("Error", {
				description: parseWailsError(error).message,
			});
			statusMessage = "Could not save update settings.";
		}
	}

	async function handleIntervalChange(value: string) {
		const checkInterval = Number(value);

		statusMessage = "";

		try {
			const updates = await UpdateService.GetUpdatesConfig();
			await UpdateService.SetUpdatesConfig({
				...updates,
				checkIntervalMinutes: checkInterval,
			});
			configStore.checkIntervalValue = checkInterval;
			statusMessage = "Check interval updated.";
		} catch (error) {
			toast.error("Error", {
				description: parseWailsError(error).message,
			});
			statusMessage = "Could not save update interval.";
		}
	}

	async function handleCheckOnStartupChange(value: boolean) {
		statusMessage = "";

		try {
			const updates = await UpdateService.GetUpdatesConfig();
			await UpdateService.SetUpdatesConfig({
				...updates,
				checkOnStartup: value,
			});
			statusMessage = value ? "Check on startup enabled." : "Check on startup disabled.";
		} catch (error) {
			toast.error("Error", {
				description: parseWailsError(error).message,
			});
			statusMessage = "Could not save check on startup setting.";
		}
	}

	async function handleChannelChange(value: string) {
		const isPrerelease = Boolean(value);
		statusMessage = "Restarting application...";

		try {
			const updates = await UpdateService.GetUpdatesConfig();
			await UpdateService.SetUpdatesConfig({
				...updates,
				includePrereleases: isPrerelease,
			});
			configStore.includePrereleases = isPrerelease;
		} catch (error) {
			toast.error("Error", {
				description: parseWailsError(error).message,
			});
			statusMessage = "Could not change update channel.";
		}
	}

	async function handleManualCheck() {
		statusMessage = "Checking for updates...";

		try {
			await UpdateService.Check();
			statusMessage = "Update check complete.";
		} catch (error) {
			toast.error("Error", {
				description: parseWailsError(error).message,
			});
			statusMessage = "Update check failed.";
		}
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
				<Button onclick={handleManualCheck} variant="outline" size="sm" class="h-7 gap-1.5 text-xs">
					<RefreshCw class="size-3.5" />
					<span>Check now</span>
				</Button>
			</div>
		</div>

		<div class="flex flex-col gap-4 rounded-lg border bg-card p-4">
			<div class="flex items-center justify-between gap-4">
				<Label for="check-on-startup">Check for updates on startup</Label>
				<Switch
					id="check-on-startup"
					bind:checked={configStore.checkOnStartup}
					onCheckedChange={handleCheckOnStartupChange}
				/>
			</div>

			<Separator />

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
					<Label>Check every</Label>
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
					<Badge
						variant="secondary"
						class="border-0 bg-amber-500/10 text-[10px] font-semibold text-amber-600 uppercase hover:bg-amber-500/10 dark:text-amber-400"
					>
						requires restart
					</Badge>
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
