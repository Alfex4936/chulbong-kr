import getReportsForMarker from "@/api/report/getReportsForMarker";
import { useQuery } from "@tanstack/react-query";

const useMarkerReportsData = (markerId: number) => {
  return useQuery({
    queryKey: ["marker", "report", markerId],
    queryFn: () => {
      return getReportsForMarker(markerId);
    },
  });
};

export default useMarkerReportsData;
