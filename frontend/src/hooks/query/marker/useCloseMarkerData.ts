import { useInfiniteQuery } from "@tanstack/react-query";
import getCloseMarker from "@/api/markers/getCloseMarker";

interface Props {
  lat: number;
  lng: number;
  distance: number;
}

const useCloseMarkerData = ({ lat, lng, distance }: Props) => {
  return useInfiniteQuery({
    queryKey: ["closeMarker", distance],
    queryFn: async ({ pageParam }) => {
      return getCloseMarker({ lat, lng, distance, pageParam });
    },
    initialPageParam: 1,
    getNextPageParam: (lastPage) => {
      if (lastPage.currentPage < lastPage.totalPages)
        return lastPage.currentPage + 1;
      return undefined;
    },
    enabled: false && lat && lng,
  });
};

export default useCloseMarkerData;
