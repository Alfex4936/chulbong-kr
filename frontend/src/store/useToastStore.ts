import { create } from "zustand";

interface ToastState {
  isToast: boolean;
  close: VoidFunction;
  open: VoidFunction;
}

const useToastStore = create<ToastState>()((set) => ({
  isToast: false,
  close: () => set({ isToast: false }),
  open: () => set({ isToast: true }),
}));

export default useToastStore;
