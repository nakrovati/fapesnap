import { A } from "@solidjs/router";
import { buttonVariants } from "~/components/ui/button";
import { cn } from "~/lib/utils";

function Nav() {
  return (
    <div class="p-4 flex flex-col min-h-[100dvh]">
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
    </div>
  );
}

export default Nav;
