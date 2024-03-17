import { useMutation, useQueryClient } from "@tanstack/react-query";
import type { Marker } from "../../../types/Marker.types";
import deleteFavorites from "../../../api/favorite/deleteFavorites";

const useDeleteFavorite = (id: number) => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: () => {
      return deleteFavorites(id);
    },
    onMutate: async () => {
      await queryClient.cancelQueries({ queryKey: ["marker", id] });

      const previousMarkerData: Marker = queryClient.getQueryData([
        "marker",
        id,
      ]) as Marker;

      if (previousMarkerData.favCount) {
        queryClient.setQueryData(["marker", id], {
          ...previousMarkerData,
          favorited: false,
          favCount: previousMarkerData.favCount - 1,
        });
      } else {
        queryClient.setQueryData(["marker", id], {
          ...previousMarkerData,
          favorited: false,
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
      queryClient.invalidateQueries({ queryKey: ["favorite"] });
    },
  });
};

export default useDeleteFavorite;
