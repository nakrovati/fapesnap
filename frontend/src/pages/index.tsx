import { StopTask } from "$wails/go/main/App";
import { Show } from "solid-js";
import { Button } from "~/components/ui/button";
import {
	Select,
	SelectContent,
	SelectItem,
	SelectTrigger,
	SelectValue,
} from "~/components/ui/select";
import { TextField, TextFieldInput } from "~/components/ui/text-field";
import { PhotoGallery, usePhotosDownloader } from "~/features/photo-downloader";
import {
	ProviderSelector,
	useProviderSelector,
} from "~/features/provider-selector";

function IndexPage() {
	const { provider, providers, setProvider } = useProviderSelector();
	const {
		setCollection,
		photos,
		maxParallelDownloads,
		setMaxParallelDownloads,
		downloading,
		loading,
		previewPhotos,
		downloadPhotos,
	} = usePhotosDownloader();

	const collectionTextFieldPlaceholder = () =>
		provider().type === "id"
			? "Enter the album ID or URL"
			: "Enter the user's name or profile URL";

	function handleDownloadPhotos() {
		downloadPhotos(provider().value);
	}

	function handlePreviewPhotos() {
		previewPhotos(provider().value);
	}

	function handleCancelTask() {
		StopTask();
	}

	return (
		<div class="flex flex-col grow h-full">
			<div class="flex gap-2">
				<TextField class="grow">
					<TextFieldInput
						placeholder={collectionTextFieldPlaceholder()}
						type="text"
						autocorrect="off"
						id="collection"
						onInput={(e) => setCollection(e.currentTarget.value)}
					/>
				</TextField>

				<ProviderSelector
					providers={providers}
					provider={provider()}
					onChange={setProvider}
				/>
			</div>

			<div class="mt-4 flex flex-col mx-auto items-center">
				<Select
					placeholder="Select the number of parallel downloads"
					value={maxParallelDownloads()}
					class="w-fit "
					onChange={setMaxParallelDownloads}
					options={[1, 2, 3, 4, 5]}
					itemComponent={(props) => (
						<SelectItem item={props.item}>{props.item.rawValue}</SelectItem>
					)}
				>
					<SelectTrigger aria-label="Max parallel downloads">
						<SelectValue<string>>
							{(state) => state.selectedOption()}
						</SelectValue>
					</SelectTrigger>
					<SelectContent />
				</Select>
				<p class="pt-2 text-sm text-gray-500">
					Number of photos uploaded simultaneously
				</p>
			</div>

			<div class="w-full flex justify-center gap-2 mt-4">
				<Button onClick={handleDownloadPhotos}>Download</Button>
				<Button onClick={handlePreviewPhotos} variant="secondary">
					Preview
				</Button>
				<Show when={downloading()}>
					<Button onClick={handleCancelTask} variant="destructive">
						Cancel
					</Button>
				</Show>
			</div>

			<div class="overflow-y-auto mt-4">
				<PhotoGallery photos={photos()} loading={loading()} />
			</div>
		</div>
	);
}

export default IndexPage;
