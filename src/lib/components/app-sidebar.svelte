<script lang="ts">
	import { resolve } from "$app/paths";
	import DownloadIcon from "@lucide/svelte/icons/download";
	import SettingsIcon from "@lucide/svelte/icons/settings";
	import { toast } from "svelte-sonner";

	import { KoFi } from "#lib/assets/icons/index.ts";
	import * as Sidebar from "#lib/components/ui/sidebar/index.ts";
	import { parseError } from "#lib/utils.ts";
	import { AppService } from "#mygo";

	const items = [
		{ title: "Download media", url: "/", icon: DownloadIcon },
		{ title: "Settings", url: "/settings", icon: SettingsIcon },
	] as const;

	async function openSupportPage() {
		try {
			await AppService.openExternal("https://ko-fi.com/Y8Y3147KNB");
		} catch (error) {
			toast.error("Error", { description: parseError(error).message });
		}
	}
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
						<button {...props} onclick={openSupportPage}>
							<KoFi />Support
						</button>
					{/snippet}
				</Sidebar.MenuButton>
			</Sidebar.MenuItem>
		</Sidebar.Menu>
	</Sidebar.Footer>
</Sidebar.Root>
