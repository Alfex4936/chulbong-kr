import { type MarkerClusterer } from "@/types/Cluster.types";
import type { KakaoMap, KakaoMarker } from "@/types/KakaoMap.types";
import { create } from "zustand";
import { devtools } from "zustand/middleware";

interface MapState {
  map: KakaoMap | null;
  clusterer: MarkerClusterer | null;
  markers: KakaoMarker[] | null;
  setMap: (map: KakaoMap) => void;
  setClusterer: (clusterer: MarkerClusterer) => void;
  setMarkers: (markers: KakaoMarker[]) => void;
}

const useMapStore = create(
  devtools<MapState>((set) => ({
    map: null,
    clusterer: null,
    markers: null,
    setMap: (map: KakaoMap) => set({ map }),
    setClusterer: (clusterer: MarkerClusterer) => set({ clusterer }),
    setMarkers: (markers: KakaoMarker[]) => set({ markers }),
  }))
);

export default useMapStore;
