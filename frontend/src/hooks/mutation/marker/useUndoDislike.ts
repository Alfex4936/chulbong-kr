import { useMutation, useQueryClient } from "@tanstack/react-query";
import markerUnDislike from "../../../api/markers/markerUnDislike";
import type { Marker } from "../../../types/Marker.types";

const useUndoDislike = (id: number) => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: () => {
      return markerUnDislike(id);
    },
    onMutate: async () => {
      await queryClient.cancelQueries({ queryKey: ["marker", id] });

      const previousMarkerData: Marker = queryClient.getQueryData([
        "marker",
        id,
      ]) as Marker;

      if (
        previousMarkerData.dislikeCount &&
        previousMarkerData.dislikeCount > 0
      ) {
        queryClient.setQueryData(["marker", id], {
          ...previousMarkerData,
          disliked: false,
          dislikeCount: previousMarkerData.dislikeCount - 1,
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

export default useUndoDislike;
