<script lang="ts">
	import { resolve } from "$app/paths";
	import { KoFi } from "$lib/assets/icons";
	import * as Sidebar from "$lib/components/ui/sidebar";
	import DownloadIcon from "@lucide/svelte/icons/download";
	import SettingsIcon from "@lucide/svelte/icons/settings";
	import { Browser } from "@wailsio/runtime";

	const items = [
		{ title: "Download media", url: "/", icon: DownloadIcon },
		{ title: "Settings", url: "/settings", icon: SettingsIcon },
	] as const;
</script>

<Sidebar.Root>
	<Sidebar.Content>
		<Sidebar.Group>
			<Sidebar.GroupContent>
				<Sidebar.Menu>
					{#each items as item (item.title)}
						<Sidebar.MenuItem>
							<Sidebar.MenuButton>
								{#snippet child({ props })}
									<a href={resolve(item.url)} {...props}>
										<item.icon />
										<span>{item.title}</span>
									</a>
								{/snippet}
							</Sidebar.MenuButton>
						</Sidebar.MenuItem>
					{/each}
				</Sidebar.Menu>
			</Sidebar.GroupContent>
		</Sidebar.Group>
	</Sidebar.Content>
	<Sidebar.Footer>
		<Sidebar.Menu>
			<Sidebar.MenuItem>
				<Sidebar.MenuButton>
					{#snippet child({ props })}
						<button {...props} onclick={() => Browser.OpenURL("https://ko-fi.com/Y8Y3147KNB")}>
							<KoFi />Support
						</button>
					{/snippet}
				</Sidebar.MenuButton>
			</Sidebar.MenuItem>
		</Sidebar.Menu>
	</Sidebar.Footer>
</Sidebar.Root>
