import { useMutation } from "@tanstack/react-query";
import updateDescription from "../../../api/markers/updateDescription";
import { useQueryClient } from "@tanstack/react-query";
import { Marker } from "../../../types/Marker.types";

const useUpdateDesc = (desc: string, id: number) => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: () => {
      return updateDescription(desc, id);
    },

    onMutate: async () => {
      await queryClient.cancelQueries({ queryKey: ["marker", id] });

      const previousMarkerData: Marker = queryClient.getQueryData([
        "marker",
        id,
      ]) as Marker;

      queryClient.setQueryData(["marker", id], {
        ...previousMarkerData,
        description: desc,
      });

      return { previousMarkerData };
    },

    onError(_error, _hero, context?: { previousMarkerData: Marker }) {
      if (context?.previousMarkerData) {
        queryClient.setQueryData(["marker", id], context.previousMarkerData);
      }
    },

    onSettled: () => {
      queryClient.invalidateQueries({ queryKey: ["marker", id] });
    },
  });
};

export default useUpdateDesc;
