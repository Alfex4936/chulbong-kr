import { useQuery } from "@tanstack/react-query";
import markerAriaRanking from "@/api/markers/markerAriaRanking";

const useMarkerAriaRankingData = (lat: number, lng: number) => {
  return useQuery({
    queryKey: ["marker-ranking", "aria"],
    queryFn: () => {
      return markerAriaRanking(lat, lng);
    },
    retry: false,
  });
};

export default useMarkerAriaRankingData;
