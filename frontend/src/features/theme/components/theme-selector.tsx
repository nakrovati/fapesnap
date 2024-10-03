import { type ConfigColorMode, useColorMode } from "@kobalte/core/color-mode";

import { createSignal } from "solid-js";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "~/components/ui/select";

const themes: ConfigColorMode[] = ["light", "dark", "system"];

function ThemeSelector() {
  const { setColorMode } = useColorMode();
  const [theme, setTheme] = createSignal<ConfigColorMode>(
    localStorage.getItem("theme") as ConfigColorMode,
  );

  return (
    <Select
      value={theme()}
      onChange={(newColorMode) => {
        setColorMode(newColorMode ?? "system");
        setTheme(newColorMode ?? "system");
      }}
      options={themes}
      placeholder="Select a theme…"
      itemComponent={(props) => (
        <SelectItem item={props.item}>{props.item.rawValue}</SelectItem>
      )}
    >
      <SelectTrigger aria-label="Fruit" class="w-48">
        <SelectValue<string>>{(state) => state.selectedOption()}</SelectValue>
      </SelectTrigger>
      <SelectContent />
    </Select>
  );
}

export { ThemeSelector };
