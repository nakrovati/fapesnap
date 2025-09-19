<script lang="ts">
	import "../app.css";
	import favicon from "$lib/assets/favicon.svg";
	import { ModeWatcher } from "mode-watcher";
	import * as Sidebar from "$lib/components/ui/sidebar";
	import AppSidebar from "$lib/components/app-sidebar.svelte";
	import * as wails from "$lib/wailsjs/runtime";
	import { Toaster } from "$lib/components/ui/sonner";
	import { toast } from "svelte-sonner";
	import { page } from "$app/state";

	let { children } = $props();

	let { title } = $derived(page.data.meta);

	$effect(() => {
		wails.EventsOn("download-start", () => {
			toast("Download started");
		});
		wails.EventsOn("download-complete", (description: string) => {
			toast.success("Download complete", {
				description,
			});
		});

		return () => {
			wails.EventsOff("download-start");
			wails.EventsOff("download-complete");
		};
	});
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
</svelte:head>

<ModeWatcher />
<Toaster />
<Sidebar.Provider>
	<AppSidebar />
	<main class="w-full p-2">
		<header class="flex items-center gap-2">
			<Sidebar.Trigger />
			<h1>{title}</h1>
		</header>
		<div class="mt-2">
			{@render children?.()}
		</div>
	</main>
</Sidebar.Provider>
