export type Providers = "fapello" | "fapodrop" | "bunkr";

export interface Provider {
  value: Providers;
  label: CapitalizeFirst<Providers>;
  type: "name" | "id";
}

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
