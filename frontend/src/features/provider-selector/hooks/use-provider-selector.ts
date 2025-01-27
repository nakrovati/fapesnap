import { createSignal } from "solid-js";
import { providers } from "../constants";

export function useProviderSelector() {
	const [provider, setProvider] = createSignal(providers[0]);

	return {
		provider,
		setProvider,
		providers,
	};
}
