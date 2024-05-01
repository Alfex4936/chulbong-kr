import markerDislike from "@/api/markers/markerDislike";
import { useToast } from "@/components/ui/use-toast";
import useLoginModalStateStore from "@/store/useLoginModalStateStore";
import type { Marker } from "@/types/Marker.types";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { isAxiosError } from "axios";

const useMarkerDislike = (id: number) => {
  const { open } = useLoginModalStateStore();
  const { toast } = useToast();
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
