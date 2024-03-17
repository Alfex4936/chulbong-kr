import type { ImageUploadState } from "../components/UploadImage/UploadImage";
import { create } from "zustand";

interface FormDataState {
  description: string;
  photoUrl: string;
  latitude: number;
  longitude: number;
  imageForm: ImageUploadState[];
  resetData: VoidFunction;
  resetPosition: VoidFunction;
  setPosition: (lat: number, lon: number) => void;
  setPhoto: (url: string) => void;
  setImageForm: (file: ImageUploadState) => void;
  replaceImages: (files: ImageUploadState[]) => void;
}

const useUploadFormDataStore = create<FormDataState>()((set) => ({
  description: "",
  photoUrl: "",
  latitude: 0,
  longitude: 0,
  imageForm: [],

  resetData: () =>
    set((state) => ({
      ...state,
      description: "",
      photoUrl: "",
      imageForm: [],
    })),
  resetPosition: () =>
    set((state) => ({
      ...state,
      latitude: 0,
      longitude: 0,
    })),

  setPosition: (lat: number, lon: number) =>
    set({ latitude: lat, longitude: lon }),

  setPhoto: (url: string) => set({ photoUrl: url }),

  setImageForm: (file: ImageUploadState) =>
    set((state) => ({ imageForm: [...state.imageForm, file] })),
  replaceImages: (files: ImageUploadState[]) => set({ imageForm: [...files] }),
}));

export default useUploadFormDataStore;
