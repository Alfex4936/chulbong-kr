import { create } from "zustand";

interface ModalState {
  loginModal: boolean;
  signupModal: boolean;
  passwordModal: boolean;
  close: VoidFunction;
  openLogin: VoidFunction;
  closeLogin: VoidFunction;
  openSignup: VoidFunction;
  closeSignup: VoidFunction;
  openPassword: VoidFunction;
  closePassword: VoidFunction;
}

const useModalStore = create<ModalState>()((set) => ({
  loginModal: false,
  signupModal: false,
  passwordModal: false,
  close: () =>
    set({ loginModal: false, signupModal: false, passwordModal: false }),
  openLogin: () => set({ loginModal: true }),
  closeLogin: () => set({ loginModal: false }),
  openSignup: () => set({ signupModal: true }),
  closeSignup: () => set({ signupModal: false }),
  openPassword: () => set({ passwordModal: true }),
  closePassword: () => set({ passwordModal: false }),
}));

export default useModalStore;
