import { create } from "zustand";
import { devtools } from "zustand/middleware";

interface ToggleState {
  isOpen: boolean;
  open: VoidFunction;
}

const useBodyToggleStore = create(
  devtools<ToggleState>((set) => ({
    isOpen: true,
    open: () => set((state) => ({ isOpen: !state.isOpen })),
  }))
);

export default useBodyToggleStore;
