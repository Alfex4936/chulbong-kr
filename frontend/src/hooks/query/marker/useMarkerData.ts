import getMarker from "@/api/markers/getMarker";
import { useQuery } from "@tanstack/react-query";

const useMarkerData = (id: number) => {
  return useQuery({
    queryKey: ["marker", id],
    queryFn: () => {
      return getMarker(id);
    },
  });
};

export default useMarkerData;
