import setFavorite from "@/api/favorite/setFavorite";
import type { Marker } from "@/types/Marker.types";
import { useMutation, useQueryClient } from "@tanstack/react-query";

const useSetFavorite = (id: number) => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: () => {
      return setFavorite(id);
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
          favorited: true,
          favCount: previousMarkerData.favCount + 1,
        });
      } else {
        queryClient.setQueryData(["marker", id], {
          ...previousMarkerData,
          favorited: true,
          favCount: 1,
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
      queryClient.invalidateQueries({ queryKey: ["marker", "bookmark"] });
    },
  });
};

export default useSetFavorite;
