import { create } from "zustand";

interface OnBoardingState {
  step: number;
  isOnBoarding: boolean;
  setStep: (step: number) => void;
  nextStep: VoidFunction;
  prevStep: VoidFunction;
}

const useOnBoardingStore = create<OnBoardingState>()((set) => ({
  step: 1,
  isOnBoarding: true,
  setStep: (step) => set({ step }),
  nextStep: () => set((state) => ({ step: state.step + 1 })),
  prevStep: () => set((state) => ({ step: state.step - 1 })),
}));

export default useOnBoardingStore;
