import { useQuery } from "@tanstack/react-query";
import markerRanking from "../../../api/markers/markerRanking";

const useMarkerRankingData = () => {
  return useQuery({
    queryKey: ["marker-ranking"],
    queryFn: markerRanking,
    retry: false,

    gcTime: 0,
  });
};

export default useMarkerRankingData;
