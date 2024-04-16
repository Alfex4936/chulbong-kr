import { create } from "zustand";
import { devtools } from "zustand/middleware";

interface TabState {
  curTab: string;
  setCurTab: (curTab: string) => void;
}

const useTabStore = create(
  devtools<TabState>((set) => ({
    curTab: "",
    setCurTab: (curTab: string) => set({ curTab }),
  }))
);

export default useTabStore;
