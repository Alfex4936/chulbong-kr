import { create } from "zustand";

interface PageLoadingState {
  isLoading: boolean;
  visible: boolean;
  setVisible: (visible: boolean) => void;
  setLoading: (loading: boolean) => void;
}

const usePageLoadingStore = create<PageLoadingState>()((set) => ({
  isLoading: false,
  visible: false,
  setLoading: (loading) => set({ isLoading: loading }),
  setVisible: (visible) => set({ visible: visible }),
}));

export default usePageLoadingStore;
