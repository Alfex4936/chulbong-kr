import { create } from "zustand";

interface BearState {
  loginModal: boolean;
  openLogin: VoidFunction;
  closeLogin: VoidFunction;
}

const useModalStore = create<BearState>()((set) => ({
  loginModal: false,
  openLogin: () => set({ loginModal: true }),
  closeLogin: () => set({ loginModal: false }),
}));

export default useModalStore;
