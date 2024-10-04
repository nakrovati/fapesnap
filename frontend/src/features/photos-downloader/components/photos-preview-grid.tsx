import { For, Show } from "solid-js";

interface PhotosPreviewProps {
  photos: string[];
  loading: boolean;
}

export function PhotosPreview(props: PhotosPreviewProps) {
  return (
    <div class="grid gap-y-4 gap-x-2 grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 2xl:grid-cols-6 ">
      <Show
        when={!props.loading}
        fallback={
          <>
            <div class="rounded animate-pulse aspect-[3/4] bg-gray-500" />
            <div class="rounded animate-pulse aspect-[3/4] bg-gray-500" />
            <div class="rounded animate-pulse aspect-[3/4] bg-gray-500 hidden md:block" />
            <div class="rounded animate-pulse aspect-[3/4] bg-gray-500 hidden lg:block" />
            <div class="rounded animate-pulse aspect-[3/4] bg-gray-500 hidden xl:block" />
            <div class="rounded animate-pulse aspect-[3/4] bg-gray-500 hidden 2xl:block" />
          </>
        }
      >
        <For each={props.photos}>
          {(item) => (
            <div>
              <img
                src={item}
                alt=""
                class="min-h-48 object-top object-contain rounded"
                loading="lazy"
              />
            </div>
          )}
        </For>
      </Show>
    </div>
  );
}
