import { create } from "zustand";
import { createJSONStorage, persist } from "zustand/middleware";

interface PWAAlertState {
  fst: boolean;
  setFstTrue: VoidFunction;
  setFstFalse: VoidFunction;
}

const usePWAAlertStore = create<PWAAlertState>()(
  persist(
    (set) => ({
      fst: false,
      setFstTrue: () => set(() => ({ fst: true })),
      setFstFalse: () => set(() => ({ fst: false })),
    }),
    {
      name: "fs-pw",
      storage: createJSONStorage(() => sessionStorage),
    }
  )
);

export default usePWAAlertStore;
