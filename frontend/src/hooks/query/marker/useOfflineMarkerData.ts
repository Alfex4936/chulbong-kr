import saveOffline from "@/api/markers/saveOffline";
import { useQuery } from "@tanstack/react-query";

const useOfflineMarkerData = (lat: number, lng: number) => {
  return useQuery({
    queryKey: ["marker", "offline", lat, lng],
    queryFn: () => {
      return saveOffline(lat, lng);
    },
  });
};

export default useOfflineMarkerData;
