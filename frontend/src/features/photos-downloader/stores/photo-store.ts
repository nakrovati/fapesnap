import { createStore } from "solid-js/store";

interface PhotosStore {
  photos: string[];
}

const [photoStore, setPhotoStore] = createStore<PhotosStore>({
  photos: [],
});

export { photoStore, setPhotoStore };
