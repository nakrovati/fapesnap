import type { ParentComponent } from "solid-js";
import { Nav } from "~/components/nav";
import { Separator } from "~/components/ui/separator";

export const DefaultLayout: ParentComponent = (props) => {
	return (
		<div class="h-screen flex">
			<Nav />

			<Separator orientation="vertical" />

			<main class="p-4 w-full">{props.children}</main>
		</div>
	);
};
