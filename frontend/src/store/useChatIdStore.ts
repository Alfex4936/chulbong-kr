import { v4 as uuidv4 } from "uuid";
import { create } from "zustand";
import { createJSONStorage, persist } from "zustand/middleware";

interface CurrentMarkerState {
  cid: string;
  setId: VoidFunction;
}

const useChatIdStore = create<CurrentMarkerState>()(
  persist(
    (set) => ({
      cid: "",
      setId: () => set({ cid: uuidv4() }),
    }),
    {
      name: "cid",
      storage: createJSONStorage(() => localStorage),
    }
  )
);

export default useChatIdStore;
