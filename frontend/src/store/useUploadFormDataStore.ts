import { create } from "zustand";

interface FormDataState {
  description: string;
  photoUrl: string;
  latitude: number;
  longitude: number;
  imageForm: File[];
  resetData: VoidFunction;
  resetPosition: VoidFunction;
  setPosition: (lat: number, lon: number) => void;
  setPhoto: (url: string) => void;
  setImageForm: (file: File) => void;
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

  setImageForm: (file: File) =>
    set((state) => ({ imageForm: [...state.imageForm, file] })),
}));

export default useUploadFormDataStore;
