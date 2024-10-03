import type { Setter } from "solid-js";
import {
	Select,
	SelectContent,
	SelectItem,
	SelectTrigger,
	SelectValue,
} from "~/components/ui/select";
import type { Provider } from "../constants";

interface ProviderSelectorProps {
	provider: Provider;
	providers: Provider[];
	onChange: Setter<Provider>;
}

export function ProvidersSelector(props: ProviderSelectorProps) {
	return (
		<Select
			value={props.provider}
			onChange={props.onChange}
			optionValue="value"
			optionTextValue="label"
			options={props.providers}
			placeholder="Select a provider..."
			itemComponent={(props) => (
				<SelectItem item={props.item}>{props.item.rawValue.value}</SelectItem>
			)}
		>
			<SelectTrigger aria-label="Provider">
				<SelectValue<Provider>>
					{(state) => state.selectedOption().label}
				</SelectValue>
			</SelectTrigger>
			<SelectContent />
		</Select>
	);
}
