import getMarkerRanking from "@/api/markers/getMarkerRanking";
import { useQuery } from "@tanstack/react-query";

const useMarkerRankingData = () => {
  return useQuery({
    queryKey: ["ranking", "top10"],
    queryFn: getMarkerRanking,
    staleTime: 0,
  });
};

export default useMarkerRankingData;
