import { create } from "zustand";
import { createJSONStorage, persist } from "zustand/middleware";

interface PsitionState {
  lat: number;
  lng: number;
  level: number;
  setPosition: (lat: number, lng: number) => void;
  setLevel: (level: number) => void;
}

const useMapStatusStore = create<PsitionState>()(
  persist(
    (set) => ({
      lat: 37.566535,
      lng: 126.9779692,
      level: 3,
      setPosition: (lat: number, lng: number) => set({ lat, lng }),
      setLevel: (level: number) => set({ level }),
    }),
    {
      name: "ps",
      storage: createJSONStorage(() => sessionStorage),
    }
  )
);

export default useMapStatusStore;
