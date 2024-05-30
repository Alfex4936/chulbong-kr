import { create } from "zustand";

interface Facilities {
  철봉: number;
  평행봉: number;
}

export interface ImageUploadState {
  file: File | null;
  previewURL: string | null;
  id: string | null;
}

interface FormDataState {
  description: string;
  photoUrl: string;
  latitude: number;
  longitude: number;
  imageForm: ImageUploadState[];
  facilities: Facilities;
  resetData: VoidFunction;
  resetPosition: VoidFunction;
  setPosition: (lat: number, lon: number) => void;
  setPhoto: (url: string) => void;
  setImageForm: (file: ImageUploadState) => void;
  replaceImages: (files: ImageUploadState[]) => void;
  increaseChulbong: () => void;
  decreaseChulbong: () => void;
  increasePenghang: () => void;
  decreasePenghang: () => void;
}

const useUploadFormDataStore = create<FormDataState>()((set) => ({
  description: "",
  photoUrl: "",
  latitude: 0,
  longitude: 0,
  imageForm: [],
  facilities: { 철봉: 0, 평행봉: 0 },
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
  increaseChulbong: () =>
    set((state) => ({
      facilities: {
        철봉: state.facilities.철봉 + 1,
        평행봉: state.facilities.평행봉,
      },
    })),
  decreaseChulbong: () =>
    set((state) => ({
      facilities: {
        철봉: state.facilities.철봉 - 1,
        평행봉: state.facilities.평행봉,
      },
    })),
  increasePenghang: () =>
    set((state) => ({
      facilities: {
        철봉: state.facilities.철봉,
        평행봉: state.facilities.평행봉 + 1,
      },
    })),
  decreasePenghang: () =>
    set((state) => ({
      facilities: {
        철봉: state.facilities.철봉,
        평행봉: state.facilities.평행봉 - 1,
      },
    })),
}));

export default useUploadFormDataStore;
