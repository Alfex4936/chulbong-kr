import setNewMarker from "@/api/markers/setNewMarker";
import { useMutation } from "@tanstack/react-query";

const useUploadMarker = () => {
  return useMutation({
    mutationFn: setNewMarker,
  });
};

export default useUploadMarker;
