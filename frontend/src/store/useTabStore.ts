import { create } from "zustand";

interface TabState {
  curTab: string;
  disableIndex: number | null;
  setCurTab: (curTab: string) => void;
  setDisable: (disableIndex: number) => void;
}

const useTabStore = create<TabState>((set) => ({
  curTab: "",
  disableIndex: null,
  setCurTab: (curTab: string) => set({ curTab }),
  setDisable: (disableIndex: number) => set({ disableIndex }),
}));

export default useTabStore;
