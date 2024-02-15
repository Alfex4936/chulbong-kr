import { create } from "zustand";

interface ModalState {
  loginModal: boolean;
  signupModal: boolean;
  close: VoidFunction;
  openLogin: VoidFunction;
  closeLogin: VoidFunction;
  openSignup: VoidFunction;
  closeSignup: VoidFunction;
}

const useModalStore = create<ModalState>()((set) => ({
  loginModal: false,
  signupModal: false,
  close: () => set({ loginModal: false, signupModal: false }),
  openLogin: () => set({ loginModal: true }),
  closeLogin: () => set({ loginModal: false }),
  openSignup: () => set({ signupModal: true }),
  closeSignup: () => set({ signupModal: false }),
}));

export default useModalStore;
