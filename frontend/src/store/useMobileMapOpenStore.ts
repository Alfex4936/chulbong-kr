import { create } from "zustand";
import { devtools } from "zustand/middleware";

interface MobileMapState {
  isOpen: boolean;
  open: VoidFunction;
  close: VoidFunction;
  toggle: VoidFunction;
}

const useMobileMapOpenStore = create(
  devtools<MobileMapState>((set) => ({
    isOpen: false,
    open: () => set({ isOpen: true }),
    close: () => set({ isOpen: false }),
    toggle: () => set((state) => ({ isOpen: !state.isOpen })),
  }))
);

export default useMobileMapOpenStore;
