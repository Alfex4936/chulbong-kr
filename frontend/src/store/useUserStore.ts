import { create } from "zustand";
import { createJSONStorage, persist } from "zustand/middleware";
import type { LoginRes } from "@/api/auth/login";

interface Props {
  user: LoginRes;
  setUser: (user: LoginRes) => void;
  resetUser: VoidFunction;
}

const useUserStore = create<Props>()(
  persist(
    (set) => ({
      user: {
        token: "",
        user: {
          email: "",
          userId: -1,
          username: "",
        },
      },
      setUser: (user: LoginRes) => set({ user }),
      resetUser: () =>
        set({
          user: {
            token: "",
            user: {
              email: "",
              userId: -1,
              username: "",
            },
          },
        }),
    }),
    {
      name: "user",
      storage: createJSONStorage(() => localStorage),
    }
  )
);

export default useUserStore;
