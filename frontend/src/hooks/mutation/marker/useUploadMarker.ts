import { useMutation } from "@tanstack/react-query";
import setNewMarker from "../../../api/markers/setNewMarker";

const useUploadMarker = () => {
  return useMutation({
    mutationFn: setNewMarker,
  });
};

export default useUploadMarker;
