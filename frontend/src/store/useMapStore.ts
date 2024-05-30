import { type MarkerClusterer } from "@/types/Cluster.types";
import { type CustomOverlay } from "@/types/CustomOverlay.types";
import type { KakaoMap, KakaoMarker } from "@/types/KakaoMap.types";
import { create } from "zustand";

interface MapState {
  map: KakaoMap | null;
  clusterer: MarkerClusterer | null;
  markers: KakaoMarker[] | null;
  overlay: any;
  myLocateOverlay: CustomOverlay | null;
  setMap: (map: KakaoMap) => void;
  setClusterer: (clusterer: MarkerClusterer) => void;
  setMarkers: (markers: KakaoMarker[]) => void;
  setOverlay: (overlay: any) => void;
  setMyLocateOverlay: (myLocateOverlay: CustomOverlay) => void;
}

const useMapStore = create<MapState>()((set) => ({
  map: null,
  clusterer: null,
  markers: null,
  overlay: null,
  myLocateOverlay: null,
  setMap: (map: KakaoMap) => set({ map }),
  setClusterer: (clusterer: MarkerClusterer) => set({ clusterer }),
  setMarkers: (markers: KakaoMarker[]) => set({ markers }),
  setOverlay: (overlay: any) => set({ overlay }),
  setMyLocateOverlay: (myLocateOverlay: CustomOverlay) =>
    set({ myLocateOverlay }),
}));

export default useMapStore;
