import { sveltekit } from "@sveltejs/kit/vite";
import tailwindcss from "@tailwindcss/vite";
import wails from "@wailsio/runtime/plugins/vite";
import { defineConfig, searchForWorkspaceRoot } from "vite";

export default defineConfig({
	plugins: [tailwindcss(), sveltekit(), wails("./bindings")],
	server: {
		host: "127.0.0.1",
		port: Number(process.env.WAILS_VITE_PORT) || 9245,
		strictPort: true,
		fs: {
			allow: [
				// search up for workspace root
				searchForWorkspaceRoot(process.cwd()),
				// your custom rules
				"./bindings/*",
			],
		},
	},
});
