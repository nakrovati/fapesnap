import { A } from "@solidjs/router";
import { buttonVariants } from "~/components/ui/button";
import { cn } from "~/lib/utils";
import { KoFiButton } from "./kofi-button";
import { Separator } from "./ui/separator";

export function Nav() {
	return (
		<nav class="p-4 flex flex-col">
			<A
				activeClass="bg-accent"
				href="/"
				end
				class={cn(buttonVariants({ variant: "ghost" }), "justify-start")}
			>
				Downloader
			</A>
			<A
				href="/settings"
				activeClass="bg-accent"
				class={cn(buttonVariants({ variant: "ghost" }), "justify-start")}
			>
				Settings
			</A>

			<Separator class="my-2" />

			<KoFiButton variant={"ghost"} class="justify-start" />
		</nav>
	);
}
