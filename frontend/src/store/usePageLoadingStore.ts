import { create } from "zustand";

interface PageLoadingState {
  isLoading: boolean;
  setLoading: (loading: boolean) => void;
}

const usePageLoadingStore = create<PageLoadingState>()((set) => ({
  isLoading: false,
  setLoading: (loading) => set({ isLoading: loading }),
}));

export default usePageLoadingStore;
