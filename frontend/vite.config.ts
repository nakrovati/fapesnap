import path from "node:path";
import { defineConfig } from "vite";
import solid from "vite-plugin-solid";

// https://vitejs.dev/config/
export default defineConfig({
	plugins: [solid()],
	resolve: {
		alias: {
			"~": path.resolve(__dirname, "./src"),
			$wails: path.resolve(__dirname, "./wailsjs"),
		},
	},
});
