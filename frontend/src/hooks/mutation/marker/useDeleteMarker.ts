import deleteMarker from "@/api/markers/deleteMarker";
import { useMutation, useQueryClient } from "@tanstack/react-query";

interface Marker {
  address: string;
  description: string;
  latitude: number;
  longitude: number;
  markerId: number;
}

interface Page {
  currentPage: number;
  markers: Marker[];
  totalMarkers: number;
  totalPages: number;
}

interface Data {
  pages: Page[];
  pageParams: number[];
}

const useDeleteMarker = (id: number) => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: () => {
      return deleteMarker(id);
    },

    onMutate: async () => {
      await queryClient.cancelQueries({ queryKey: ["myMarker"] });

      const previousMarkerData: Data = queryClient.getQueryData([
        "myMarker",
      ]) as Data;

      const copy = { ...previousMarkerData };

      copy.pages.forEach((page) => {
        page.markers = page.markers.filter((marker) => marker.markerId !== id);
      });

      queryClient.setQueryData(["myMarker"], copy);

      return { previousMarkerData };
    },

    onError: (_error, _hero, context?: { previousMarkerData: Data }) => {
      if (context?.previousMarkerData) {
        queryClient.setQueryData(["myMarker"], context.previousMarkerData);
      }
    },

    onSettled: () => {
      queryClient.invalidateQueries({ queryKey: ["myMarker"] });
    },
  });
};

export default useDeleteMarker;
