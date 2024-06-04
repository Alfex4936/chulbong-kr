import { create } from "zustand";

interface Marker {
  markerId: number;
  lat: number;
  lng: number;
}

interface SelectedState {
  marker: Marker | null;
  setMarker: (marker: Marker | null) => void;
}

const useSelectedMarkerStore = create<SelectedState>((set) => ({
  marker: null,
  setMarker: (marker: Marker | null) => set({ marker }),
}));

export default useSelectedMarkerStore;
