import { Route, Router } from "@solidjs/router";
/* @refresh reload */
import { render } from "solid-js/web";
import { Toaster } from "~/components/ui/toast";
import "./style.css";
import IndexPage from "~/pages/index-page";
import DefaultLayout from "./default-layout";

const root = document.getElementById("root");
if (!root) {
	throw new Error("Root element not found");
}

render(
	() => (
		<>
			<Router root={DefaultLayout}>
				<Route path="/" component={IndexPage} />
			</Router>
			<Toaster />
		</>
	),
	root,
);
