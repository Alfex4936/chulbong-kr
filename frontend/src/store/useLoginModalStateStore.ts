import { create } from "zustand";
import { devtools } from "zustand/middleware";

interface ModalState {
  isOpen: boolean;
  open: VoidFunction;
  close: VoidFunction;
}

const useLoginModalStateStore = create(
  devtools<ModalState>((set) => ({
    isOpen: false,
    open: () => set({ isOpen: true }),
    close: () => set({ isOpen: false }),
  }))
);

export default useLoginModalStateStore;
