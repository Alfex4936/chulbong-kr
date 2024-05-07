import { type MarkerClusterer } from "@/types/Cluster.types";
import type { KakaoMap, KakaoMarker } from "@/types/KakaoMap.types";
import { create } from "zustand";
import { devtools } from "zustand/middleware";

interface MapState {
  map: KakaoMap | null;
  clusterer: MarkerClusterer | null;
  markers: KakaoMarker[] | null;
  overlay: any;
  setMap: (map: KakaoMap) => void;
  setClusterer: (clusterer: MarkerClusterer) => void;
  setMarkers: (markers: KakaoMarker[]) => void;
  setOverlay: (overlay: any) => void;
}

const useMapStore = create(
  devtools<MapState>((set) => ({
    map: null,
    clusterer: null,
    markers: null,
    overlay: null,
    setMap: (map: KakaoMap) => set({ map }),
    setClusterer: (clusterer: MarkerClusterer) => set({ clusterer }),
    setMarkers: (markers: KakaoMarker[]) => set({ markers }),
    setOverlay: (overlay: any) => set({ overlay }),
  }))
);

export default useMapStore;
