<script lang="ts">
	import "./layout.css";
	import favicon from "$lib/assets/favicon.svg";
	import { ModeWatcher } from "mode-watcher";
	import * as Sidebar from "$lib/components/ui/sidebar";
	import AppSidebar from "$lib/components/app-sidebar.svelte";
	import { Toaster } from "$lib/components/ui/sonner";
	import { toast } from "svelte-sonner";
	import { page } from "$app/state";
	import { Events } from "@wailsio/runtime";

	let { children } = $props();

	let { title } = $derived(page.data.meta);

	$effect(() => {
		Events.On("download-start", () => {
			toast("Download started");
		});
		Events.On("download-complete", (data) => {
			console.log(data);
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
