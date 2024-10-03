import { createStore } from "solid-js/store";

interface PhotosStore {
  photos: string[];
}

const [store, setStore] = createStore<PhotosStore>({
  photos: [],
});

export { store as photoStore, setStore as setPhotoStore };
