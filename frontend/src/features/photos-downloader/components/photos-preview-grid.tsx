import { For, Show } from "solid-js";

interface PhotosPreviewProps {
	photos: string[];
	loading: boolean;
}

function PhotosPreview(props: PhotosPreviewProps) {
	return (
		<div class="grid gap-y-4 mt-4 gap-x-2 grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6">
			<Show
				when={!props.loading}
				fallback={
					<>
						<div class="rounded animate-pulse aspect-[3/4] bg-gray-500" />
						<div class="rounded animate-pulse aspect-[3/4] bg-gray-500" />
						<div class="rounded animate-pulse aspect-[3/4] bg-gray-500 hidden sm:block" />
						<div class="rounded animate-pulse aspect-[3/4] bg-gray-500 hidden md:block" />
						<div class="rounded animate-pulse aspect-[3/4] bg-gray-500 hidden lg:block" />
						<div class="rounded animate-pulse aspect-[3/4] bg-gray-500 hidden xl:block" />
					</>
				}
			>
				<For each={props.photos}>
					{(item) => (
						<div>
							<img
								src={item}
								alt=""
								class="aspect-[3/4] object-contain object-top rounded"
								loading="lazy"
							/>
						</div>
					)}
				</For>
			</Show>
		</div>
	);
}

export { PhotosPreview };
