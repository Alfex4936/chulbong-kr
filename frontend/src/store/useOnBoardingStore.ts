import { create } from "zustand";

interface OnBoardingState {
  step: number;
  isOnBoarding: boolean;
  open: VoidFunction;
  close: VoidFunction;
  setStep: (step: number) => void;
  nextStep: VoidFunction;
  prevStep: VoidFunction;
}

const useOnBoardingStore = create<OnBoardingState>()((set) => ({
  step: 0,
  isOnBoarding: false,
  open: () => set(() => ({ step: 0, isOnBoarding: true })),
  close: () => set(() => ({ step: 0, isOnBoarding: false })),
  setStep: (step) => set({ step }),
  nextStep: () => set((state) => ({ step: state.step + 1 })),
  prevStep: () => set((state) => ({ step: state.step - 1 })),
}));

export default useOnBoardingStore;
