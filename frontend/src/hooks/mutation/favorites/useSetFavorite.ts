import setFavorite from "@/api/favorite/setFavorite";
import { useToast } from "@/components/ui/use-toast";
import useLoginModalStateStore from "@/store/useLoginModalStateStore";
import type { Marker } from "@/types/Marker.types";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { isAxiosError } from "axios";

const useSetFavorite = (id: number) => {
  const { open } = useLoginModalStateStore();
  const { toast } = useToast();

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

    onError(error, _hero, context?: { previousMarkerData: Marker }) {
      if (isAxiosError(error)) {
        if (error.response?.status === 401) {
          open();
        } else {
          toast({ description: "잠시 후 다시 시도해 주세요." });
        }
      } else {
        toast({ description: "잠시 후 다시 시도해 주세요." });
      }
      // open();
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
