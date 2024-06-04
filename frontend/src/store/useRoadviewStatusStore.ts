import { create } from "zustand";

interface RoadviewState {
  lat: number;
  lng: number;
  isOpen: boolean;
  setPosition: (lat: number, ing: number) => void;
  open: VoidFunction;
  close: VoidFunction;
}

const useRoadviewStatusStore = create<RoadviewState>((set) => ({
  lat: 37.566535,
  lng: 126.9779692,
  isOpen: false,
  setPosition: (lat: number, lng: number) =>
    set((state) => ({ ...state, lat, lng })),
  open: () => set({ isOpen: true }),
  close: () => set({ isOpen: false }),
}));

export default useRoadviewStatusStore;
