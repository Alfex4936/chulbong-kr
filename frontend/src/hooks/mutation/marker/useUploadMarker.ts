import { useMutation } from "@tanstack/react-query";
import setNewMarker from "../../../api/markers/setNewMarker";
import { useQueryClient } from "@tanstack/react-query";

const useUploadMarker = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: setNewMarker,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["marker", "all"] });
    },
  });
};

export default useUploadMarker;
