import { create } from "zustand";

interface ModalState {
  loading: boolean;
  name: string;
  setLoading: (loading: boolean) => void;
  setName: (name: string) => void;
}

const useLoginLoadingStore = create<ModalState>()((set) => ({
  loading: false,
  name: "",
  setLoading: (loading) => set({ loading }),
  setName: (name) => set({ name }),
}));

export default useLoginLoadingStore;
