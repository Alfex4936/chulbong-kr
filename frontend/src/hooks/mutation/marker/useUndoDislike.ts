import { useMutation, useQueryClient } from "@tanstack/react-query";
import markerUnDislike from "../../../api/markers/markerUnDislike";
import type { Marker } from "../../../types/Marker.types";

interface Dislike {
  disliked: boolean;
}

const useUndoDislike = (id: number) => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: () => {
      return markerUnDislike(id);
    },
    onMutate: async () => {
      await queryClient.cancelQueries({ queryKey: ["dislikeState", id] });
      await queryClient.cancelQueries({ queryKey: ["marker", id] });

      const previousLikeData: Dislike = queryClient.getQueryData([
        "dislikeState",
        id,
      ]) as Dislike;
      const previousMarkerData: Marker = queryClient.getQueryData([
        "marker",
        id,
      ]) as Marker;

      queryClient.setQueryData(["dislikeState", id], { disliked: false });
      if (
        previousMarkerData.dislikeCount &&
        previousMarkerData.dislikeCount > 0
      ) {
        queryClient.setQueryData(["marker", id], {
          ...previousMarkerData,
          dislikeCount: previousMarkerData.dislikeCount - 1,
        });
      }

      return { previousLikeData, previousMarkerData };
    },

    onError(
      _error,
      _hero,
      context?: { previousLikeData: Dislike; previousMarkerData: Marker }
    ) {
      if (context?.previousLikeData) {
        queryClient.setQueryData(
          ["dislikeState", id],
          context.previousLikeData
        );
      }
      if (context?.previousMarkerData) {
        queryClient.setQueryData(["marker", id], context.previousMarkerData);
      }
    },

    onSettled() {
      queryClient.invalidateQueries({ queryKey: ["dislikeState", id] });
      queryClient.invalidateQueries({ queryKey: ["marker", id] });
    },
  });
};

export default useUndoDislike;
