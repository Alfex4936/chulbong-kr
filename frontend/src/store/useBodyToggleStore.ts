import { create } from "zustand";

interface ToggleState {
  isOpen: boolean;
  open: VoidFunction;
}

const useBodyToggleStore = create<ToggleState>((set) => ({
  isOpen: true,
  open: () => set((state) => ({ isOpen: !state.isOpen })),
}));

export default useBodyToggleStore;
