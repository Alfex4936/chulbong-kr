import { useQuery } from "@tanstack/react-query";
import getAllMarker from "@/api/markers/getAllMarker";

const useAllMarkerData = () => {
  return useQuery({
    queryKey: ["markers"],
    queryFn: getAllMarker,
  });
};

export default useAllMarkerData;
