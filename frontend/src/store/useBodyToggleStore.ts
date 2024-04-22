// import { create } from "zustand";

// interface ToggleState {
//   isOpen: boolean;
//   open: VoidFunction;
// }

// const useBodyToggleState = create<ToggleState>()((set) => ({
//   isOpen: false,
//   open: () => set((state) => ({ isOpen: !state.isOpen })),
// }));

// export default useBodyToggleState;

import { create } from "zustand";
import { devtools } from "zustand/middleware";

interface ToggleState {
  isOpen: boolean;
  open: VoidFunction;
}

const useBodyToggleStore = create(
  devtools<ToggleState>((set) => ({
    isOpen: true,
    open: () => set((state) => ({ isOpen: !state.isOpen })),
  }))
);

export default useBodyToggleStore;
