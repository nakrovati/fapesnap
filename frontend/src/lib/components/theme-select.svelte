<script lang="ts">
	import * as Select from "$lib/components/ui/select";
	import { setMode, userPrefersMode } from "mode-watcher";

	type Mode = "light" | "dark" | "system";
	let currentMode = $derived(userPrefersMode.current);
</script>

<Select.Root
	bind:value={currentMode}
	type="single"
	onValueChange={(val) => (val === "system" ? setMode("system") : setMode(val as Mode))}
>
	<Select.Trigger class="w-[180px]">
		{currentMode}
	</Select.Trigger>
	<Select.Content>
		<Select.Item value="light">Light</Select.Item>
		<Select.Item value="dark">Dark</Select.Item>
		<Select.Item value="system">System</Select.Item>
	</Select.Content>
</Select.Root>
