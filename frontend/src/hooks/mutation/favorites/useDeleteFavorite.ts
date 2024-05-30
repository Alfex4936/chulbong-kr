import deleteFavorite from "@/api/favorite/deleteFavorite";
import { useToast } from "@/components/ui/use-toast";
import type { Marker } from "@/types/Marker.types";
import { useMutation, useQueryClient } from "@tanstack/react-query";

const useDeleteFavorite = (id: number) => {
  const { toast } = useToast();
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: () => {
      return deleteFavorite(id);
    },
    onMutate: async () => {
      await queryClient.cancelQueries({ queryKey: ["marker", id] });

      const previousMarkerData: Marker = queryClient.getQueryData([
        "marker",
        id,
      ]) as Marker;

      if (!previousMarkerData) return;

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
      toast({ description: "잠시 후 다시 시도해주세요." });
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

export default useDeleteFavorite;
