declare type CapitalizeFirst<T extends string> =
	T extends `${infer First}${infer Rest}` ? `${Uppercase<First>}${Rest}` : T;
