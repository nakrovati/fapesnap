import { Route, Router } from "@solidjs/router";
import { render } from "solid-js/web";
import { Toaster, showToast } from "~/components/ui/toast";
import "./style.css";
import * as wails from "$wails/runtime";
import {
  ColorModeProvider,
  ColorModeScript,
  createLocalStorageManager,
} from "@kobalte/core";
import { createEffect, onCleanup } from "solid-js";
import IndexPage from "~/pages";
import SettingsPage from "~/pages/settings";
import { DefaultLayout } from "./default-layout";

const root = document.getElementById("root");
if (!root) {
  throw new Error("Root element not found");
}

function App() {
  const storageManager = createLocalStorageManager("theme");

  createEffect(() => {
    wails.EventsOn("download-start", () => {
      showToast({ title: "Download started" });
    });

    wails.EventsOn("download-complete", (description: string) => {
      showToast({
        title: "Download complete",
        description,
        variant: "success",
      });
    });

    wails.EventsOn("download-canceled", (description: string) => {
      showToast({
        title: "Download canceled",
        description,
        variant: "warning",
      });
    });
  });

  onCleanup(() => {
    wails.EventsOff("download-start");
    wails.EventsOff("download-complete");
    wails.EventsOff("download-canceled");
  });

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
