import { create } from "zustand";

interface FormDataState {
  description: string;
  photoUrl: string;
  latitude: number;
  longitude: number;
  imageForm: File | null;
  resetData: VoidFunction;
  setPosition: (lat: number, lon: number) => void;
  setPhoto: (url: string) => void;
  setImageForm: (file: File) => void;
}

const useUploadFormDataStore = create<FormDataState>()((set) => ({
  description: "",
  photoUrl: "",
  latitude: 0,
  longitude: 0,
  imageForm: null,

  resetData: () =>
    set({
      description: "",
      photoUrl: "",
      latitude: 0,
      longitude: 0,
      imageForm: null,
    }),

  setPosition: (lat: number, lon: number) =>
    set({ latitude: lat, longitude: lon }),

  setPhoto: (url: string) => set({ photoUrl: url }),

  setImageForm: (file: File) => set({ imageForm: file }),
}));

export default useUploadFormDataStore;
