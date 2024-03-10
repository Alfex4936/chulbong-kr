import { useMutation, useQueryClient } from "@tanstack/react-query";
import type { Marker } from "../../../types/Marker.types";
import setFavorite from "../../../api/favorite/setFavorite";

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

      queryClient.setQueryData(["marker", id], {
        ...previousMarkerData,
        favorited: true,
      });

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

export default useSetFavorite;
