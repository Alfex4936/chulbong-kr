import { useQuery } from "@tanstack/react-query";
import getMarker from "../../../api/markers/getMarker";

const useGetMarker = (id: number) => {
  return useQuery({
    queryKey: ["marker", id],
    queryFn: () => {
      return getMarker(id);
    },
    retry: false,
  });
};

export default useGetMarker;
