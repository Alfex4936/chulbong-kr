import { useMutation } from "@tanstack/react-query";
import deleteMarker from "../../../api/markers/deleteMarker";

const useDeleteMarker = (id: number) => {
  return useMutation({
    mutationFn: () => {
      return deleteMarker(id);
    },
  });
};

export default useDeleteMarker;
