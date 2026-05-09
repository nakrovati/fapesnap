export interface Provider {
	label: CapitalizeFirst<Providers>;
	type: "id" | "name";
	value: Providers;
}

export type Providers = "bunkr" | "fapello" | "fapodrop";

export const providers: Provider[] = [
	{
		value: "fapello",
		label: "Fapello",
		type: "name",
	},
	{
		value: "fapodrop",
		label: "Fapodrop",
		type: "name",
	},
	{
		value: "bunkr",
		label: "Bunkr",
		type: "id",
	},
];
