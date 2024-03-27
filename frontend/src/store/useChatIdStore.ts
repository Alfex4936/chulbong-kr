import { create } from "zustand";
import { createJSONStorage, persist } from "zustand/middleware";

interface CurrentMarkerState {
  cid: string;
  setId: (cid: string) => void;
}

const useChatIdStore = create<CurrentMarkerState>()(
  persist(
    (set) => ({
      cid: "",
      setId: (cid) => set({ cid }),
    }),
    {
      name: "cid",
      storage: createJSONStorage(() => localStorage),
    }
  )
);

export default useChatIdStore;
