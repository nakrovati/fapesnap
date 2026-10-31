import { defineConfig } from "mygo-cli";

export default defineConfig({
	name: "Fapesnap",
	slug: "fapesnap",
	identifier: "com.github.nakrovati.fapesnap",
	version: "2.4.0-rc.2",
	devUrl: "http://localhost:5173",
	devCommand: "bun run dev:web",
	buildCommand: "bun run build:web",
	frontendDist: "dist",
	bindings: "src/mygo.ts",
	out: "build",
	updates: {
		publicKey: "jshheoQTd2OA1RwWzdgt/I9E3a8vcy/JHP2inYuO0iI=",
		github: "nakrovati/fapesnap-mygo",
	},
});
