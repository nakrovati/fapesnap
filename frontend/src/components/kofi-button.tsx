import { KoFi } from "~/assets/icons";
import { cn } from "~/lib/utils";
import { Button, type ButtonProps } from "./ui/button";

export function KoFiButton(props: ButtonProps) {
	function handleOpenKoFiPage() {
		// @ts-ignore
		window.runtime.BrowserOpenURL("https://ko-fi.com/Y8Y3147KNB");
	}

	return (
		<Button
			type="button"
			class={cn("gap-1", props.class)}
			onClick={handleOpenKoFiPage}
			variant={props.variant}
		>
			<div>Support</div>
			<KoFi />
		</Button>
	);
}
