import { create } from "zustand";

interface Image {
  markerId?: number;
  photoId: number | string;
  photoUrl: string;
  uploadedAt?: Date;
}

interface MarkerImageState {
  imageView: boolean;
  images: Image[] | null;
  curImage: Image | null;
  curImageIndex: number;
  setImages: (images: Image[]) => void;
  setCurImage: (curImage: Image) => void;
  setCurImageIndex: (curImageIndex: number) => void;
  openImageModal: VoidFunction;
  closeImageModal: VoidFunction;
  nextImage: VoidFunction;
  prevImage: VoidFunction;
}

const useMarkerImageStore = create<MarkerImageState>()((set) => ({
  images: null,
  imageView: false,
  curImage: null,
  curImageIndex: -1,
  setImages: (images: Image[]) => set({ images }),
  setCurImage: (curImage: Image) => set({ curImage }),
  setCurImageIndex: (curImageIndex: number) => set({ curImageIndex }),
  openImageModal: () => set({ imageView: true }),
  closeImageModal: () => set({ imageView: false }),
  nextImage: () =>
    set((state) => {
      if (!state.images) return { images: null };

      const currentIndex = state.images.findIndex(
        (image) => image.photoId === state.curImage?.photoId
      );

      if (currentIndex !== -1) {
        const nextIndex = (currentIndex + 1) % state.images.length;
        return { curImage: state.images[nextIndex] };
      }

      return { curImage: state.images[state.images.length - 1] };
    }),
  prevImage: () =>
    set((state) => {
      if (!state.images) return { images: null };

      const currentIndex = state.images.findIndex(
        (image) => image.photoId === state.curImage?.photoId
      );

      if (currentIndex !== -1) {
        const prevIndex =
          (currentIndex - 1 + state.images.length) % state.images.length;
        return { curImage: state.images[prevIndex] };
      }

      return { curImage: state.images[0] };
    }),
}));

export default useMarkerImageStore;
