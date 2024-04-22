import { type KakaoMap } from "@/types/KakaoMap.types";
import { create } from "zustand";
import { devtools } from "zustand/middleware";
import { type MarkerClusterer } from "@/types/Cluster.types";

interface MapState {
  map: KakaoMap | null;
  clusterer: MarkerClusterer | null;
  setMap: (map: KakaoMap) => void;
  setClusterer: (clusterer: MarkerClusterer) => void;
}

const useMapStore = create(
  devtools<MapState>((set) => ({
    map: null,
    clusterer: null,
    setMap: (map: KakaoMap) => set({ map }),
    setClusterer: (clusterer: MarkerClusterer) => set({ clusterer }),
  }))
);

export default useMapStore;
