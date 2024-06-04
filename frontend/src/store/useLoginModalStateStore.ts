import { create } from "zustand";

interface ModalState {
  isOpen: boolean;
  url: string;
  open: VoidFunction;
  close: VoidFunction;
  setUrl: (url: string) => void;
}

const useLoginModalStateStore = create<ModalState>((set) => ({
  isOpen: false,
  url: "home",
  open: () => set({ isOpen: true }),
  close: () => set({ isOpen: false }),
  setUrl: (url: string) => set({ url }),
}));

export default useLoginModalStateStore;
