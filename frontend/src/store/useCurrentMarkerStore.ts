import { create } from "zustand";
import { createJSONStorage, persist } from "zustand/middleware";

interface CurrentMarkerState {
  markerId: number;
  setMarker: (markerId: number) => void;
}

const useCurrentMarkerStore = create<CurrentMarkerState>()(
  persist(
    (set) => ({
      markerId: -1,
      setMarker: (markerId) => set({ markerId }),
    }),
    {
      name: "cms",
      storage: createJSONStorage(() => sessionStorage),
    }
  )
);

export default useCurrentMarkerStore;
