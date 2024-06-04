import { KakaoMap } from "@/types/KakaoMap.types";
import { create } from "zustand";

interface PsitionState {
  map: KakaoMap | null;
  lat: number;
  lng: number;
  level: number;
  loaded: boolean;
  setLoad: VoidFunction;
  setPosition: (lat: number, lng: number) => void;
  setLevel: (level: number) => void;
  setMap: (map: KakaoMap) => void;
}

const useMiniMapStatusStore = create<PsitionState>((set) => ({
  map: null,
  lat: 37.566535,
  lng: 126.9779692,
  level: 3,
  loaded: false,
  setLoad: () => set({ loaded: true }),
  setPosition: (lat: number, lng: number) => set({ lat, lng }),
  setLevel: (level: number) => set({ level }),
  setMap: (map: KakaoMap) => set({ map }),
}));

export default useMiniMapStatusStore;
