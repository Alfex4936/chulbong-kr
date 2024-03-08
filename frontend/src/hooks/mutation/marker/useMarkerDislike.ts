import { useMutation, useQueryClient } from "@tanstack/react-query";
import markerDislike from "../../../api/markers/markerDislike";
import type { Marker } from "../../../types/Marker.types";

const useMarkerDislike = (id: number) => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: () => {
      return markerDislike(id);
    },
    onMutate: async () => {
      await queryClient.cancelQueries({ queryKey: ["marker", id] });

      const previousMarkerData: Marker = queryClient.getQueryData([
        "marker",
        id,
      ]) as Marker;

      if (previousMarkerData.dislikeCount) {
        queryClient.setQueryData(["marker", id], {
          ...previousMarkerData,
          disliked: true,
          dislikeCount: previousMarkerData.dislikeCount + 1,
        });
      } else {
        queryClient.setQueryData(["marker", id], {
          ...previousMarkerData,
          disliked: true,
          dislikeCount: 1,
        });
      }

      return { previousMarkerData };
    },

    onError(_error, _hero, context?: { previousMarkerData: Marker }) {
      if (context?.previousMarkerData) {
        queryClient.setQueryData(["marker", id], context.previousMarkerData);
      }
    },

    onSettled() {
      queryClient.invalidateQueries({ queryKey: ["marker", id] });
    },
  });
};

export default useMarkerDislike;
