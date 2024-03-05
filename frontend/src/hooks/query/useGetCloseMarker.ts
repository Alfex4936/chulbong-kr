import { useInfiniteQuery } from "@tanstack/react-query";
import getCloseMarker from "../../api/markers/getCloseMarker";

interface Props {
  lat: number;
  lon: number;
  distance: number;
}

const useGetCloseMarker = ({ lat, lon, distance }: Props) => {
  return useInfiniteQuery({
    queryKey: ["closeMarker", distance],
    queryFn: ({ pageParam }) => {
      return getCloseMarker({ lat, lon, distance, pageParam });
    },
    initialPageParam: 1,
    getNextPageParam: (lastPage) => {
      if (lastPage.currentPage < lastPage.totalPages)
        return lastPage.currentPage + 1;
      return undefined;
    },
    enabled: false,
    // staleTime: 3000,
  });
};

export default useGetCloseMarker;
