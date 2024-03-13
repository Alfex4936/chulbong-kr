import { useMutation, useQueryClient } from "@tanstack/react-query";
import deleteMarker from "../../../api/markers/deleteMarker";

const useDeleteMarker = (id: number) => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: () => {
      return deleteMarker(id);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["marker", "all"] });
    },
  });
};

export default useDeleteMarker;
