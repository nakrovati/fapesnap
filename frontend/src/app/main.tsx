import { Route, Router } from "@solidjs/router";
import { render } from "solid-js/web";
import { Toaster } from "~/components/ui/toast";
import "./style.css";
import {
	ColorModeProvider,
	ColorModeScript,
	createLocalStorageManager,
} from "@kobalte/core";
import IndexPage from "~/pages";
import SettingsPage from "~/pages/settings";
import { DefaultLayout } from "./default-layout";

const root = document.getElementById("root");
if (!root) {
	throw new Error("Root element not found");
}

function App() {
	const storageManager = createLocalStorageManager("theme");

	return (
		<>
			<Router
				root={(props) => (
					<>
						<ColorModeScript storageType={storageManager.type} />
						<ColorModeProvider storageManager={storageManager}>
							<DefaultLayout>{props.children}</DefaultLayout>
						</ColorModeProvider>
						<Toaster />
					</>
				)}
			>
				<Route path="/" component={IndexPage} />
				<Route path="/settings" component={SettingsPage} />
			</Router>
		</>
	);
}

render(App, root);
