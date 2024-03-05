import { create } from "zustand";
import { createJSONStorage, persist } from "zustand/middleware";

interface PsitionState {
  lat: number;
  lng: number;
  setPosition: (lat: number, lon: number) => void;
}

const useMapPositionStore = create<PsitionState>()(
  persist(
    (set) => ({
      lat: 0,
      lng: 0,
      setPosition: (lat: number, lng: number) => set({ lat, lng }),
    }),
    {
      name: "ps",
      storage: createJSONStorage(() => sessionStorage),
    }
  )
);

export default useMapPositionStore;
