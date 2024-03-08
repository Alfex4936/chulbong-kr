import { create } from "zustand";

interface ToastState {
  isToast: boolean;
  toastText: string;
  setToastText: (text: string) => void;
  close: VoidFunction;
  open: VoidFunction;
}

const usePaymentUserStore = create<ToastState>()((set) => ({
  isToast: false,
  toastText: "",
  setToastText: (text: string) => set({ toastText: text }),
  close: () => set({ isToast: false }),
  open: () => set({ isToast: true }),
}));

export default usePaymentUserStore;
