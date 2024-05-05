import { create } from "zustand";
import { devtools } from "zustand/middleware";

interface Marker {
  markerId: number;
  lat: number;
  lng: number;
}

interface SelectedState {
  marker: Marker | null;
  setMarker: (marker: Marker | null) => void;
}

const useSelectedMarkerStore = create(
  devtools<SelectedState>((set) => ({
    marker: null,
    setMarker: (marker: Marker | null) => set({ marker }),
  }))
);

export default useSelectedMarkerStore;
