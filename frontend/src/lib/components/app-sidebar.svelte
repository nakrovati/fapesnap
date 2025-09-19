<script lang="ts">
	import * as Sidebar from "$lib/components/ui/sidebar";
	import DownloadIcon from "@lucide/svelte/icons/download";
	import { KoFi } from "$lib/assets/icons";
	import SettingsIcon from "@lucide/svelte/icons/settings";
	import { BrowserOpenURL } from "$lib/wailsjs/runtime/runtime";

	const items = [
		{ title: "Download photos", url: "/", icon: DownloadIcon },
		{ title: "Settings", url: "/settings", icon: SettingsIcon },
	];
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
									<a href={item.url} {...props}>
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
						<button {...props} onclick={() => BrowserOpenURL("https://ko-fi.com/Y8Y3147KNB")}>
							<KoFi />Support
						</button>
					{/snippet}
				</Sidebar.MenuButton>
			</Sidebar.MenuItem>
		</Sidebar.Menu>
	</Sidebar.Footer>
</Sidebar.Root>
