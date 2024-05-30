import { create } from "zustand";

interface ButtonState {
  isActive: boolean;
  setIsActive: (isActive: boolean) => void;
}

const useScrollButtonStore = create<ButtonState>()((set) => ({
  isActive: false,
  setIsActive: (isActive: boolean) => set({ isActive }),
}));

export default useScrollButtonStore;
