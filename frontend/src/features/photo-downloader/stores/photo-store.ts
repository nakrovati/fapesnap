import { createStore } from "solid-js/store";

interface PhotoStore {
	photos: string[];
}

const [photoStore, setPhotoStore] = createStore<PhotoStore>({
	photos: [],
});

export { photoStore, setPhotoStore };
