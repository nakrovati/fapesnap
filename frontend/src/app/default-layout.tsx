import type { ParentComponent } from "solid-js";
import Nav from "~/components/Nav";
import { Separator } from "~/components/ui/separator";

const DefaultLayout: ParentComponent = (props) => {
	return (
		<div class="grid grid-cols-[auto_auto_1fr]">
			<Nav />

			<Separator orientation="vertical" />

			<main class="p-4">{props.children}</main>
		</div>
	);
};

export default DefaultLayout;
