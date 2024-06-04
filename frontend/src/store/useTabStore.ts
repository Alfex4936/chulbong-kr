import { create } from "zustand";

interface TabState {
  curTab: string;
  setCurTab: (curTab: string) => void;
}

const useTabStore = create<TabState>((set) => ({
  curTab: "",
  setCurTab: (curTab: string) => set({ curTab }),
}));

export default useTabStore;
