import { create } from "zustand";
import { createJSONStorage, persist } from "zustand/middleware";

interface P {
  h: string;
  kj: { ej: string; jj: number; ol: string };
}

interface UserState {
  ka: P;
  setUser: (user: P) => void;
  resetUser: VoidFunction;
}

const useUserStore = create<UserState>()(
  persist(
    (set) => ({
      ka: {
        h: "",
        kj: {
          ej: "",
          jj: -1,
          ol: "",
        },
      },
      setUser: (ka: P) => set({ ka }),
      resetUser: () =>
        set({
          ka: {
            h: "",
            kj: {
              ej: "",
              jj: -1,
              ol: "",
            },
          },
        }),
    }),
    {
      name: "uaui",
      storage: createJSONStorage(() => localStorage),
    }
  )
);

export default useUserStore;
