<script lang="ts">
	import "./layout.css";
	import { page } from "$app/state";
	import favicon from "$lib/assets/favicon.svg";
	import AppSidebar from "$lib/components/app-sidebar.svelte";
	import * as Sidebar from "$lib/components/ui/sidebar";
	import { Toaster } from "$lib/components/ui/sonner";
	import { Events } from "@wailsio/runtime";
	import { ModeWatcher } from "mode-watcher";
	import { toast } from "svelte-sonner";

	let { children } = $props();

	let { title } = $derived(page.data.meta);

	$effect(() => {
		Events.On("download-start", () => {
			toast("Download started");
		});
		Events.On("download-complete", (data) => {
			toast.success("Download complete", {
				description: data.data.description,
			});
		});

		return () => {
			Events.Off("download-start");
			Events.Off("download-complete");
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
			{@render children()}
		</div>
	</main>
</Sidebar.Provider>
