import getAriaMarkerRanking from "@/api/markers/getAreaMarkerRanking";
import { useQuery } from "@tanstack/react-query";

const useAreaMarkerRankingData = (lat: number, lng: number, start: boolean) => {
  return useQuery({
    queryKey: ["ranking", "aria"],
    queryFn: () => {
      return getAriaMarkerRanking(lat, lng);
    },

    staleTime: 0,
    enabled: start,
  });
};

export default useAreaMarkerRankingData;
