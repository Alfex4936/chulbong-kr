import { create } from "zustand";
import { createJSONStorage, persist } from "zustand/middleware";

interface PsitionState {
  lat: number;
  lng: number;
  level: number;
  setPosition: (lat: number, lon: number) => void;
  setLevel: (level: number) => void;
}

const useMapPositionStore = create<PsitionState>()(
  persist(
    (set) => ({
      lat: 37.566535,
      lng: 126.9779692,
      level: 5,
      setPosition: (lat: number, lng: number) => set({ lat, lng }),
      setLevel: (level: number) => set({ level }),
    }),
    {
      name: "ps",
      storage: createJSONStorage(() => sessionStorage),
    }
  )
);

export default useMapPositionStore;
