import { useInfiniteQuery } from "@tanstack/react-query";
import getMyMarker from "../../api/markers/getMyMarker";

const useGetMyMarker = () => {
  return useInfiniteQuery({
    queryKey: ["myMarker"],
    queryFn: getMyMarker,
    initialPageParam: 1,
    getNextPageParam: (lastPage) => {
      if (lastPage.currentPage < lastPage.totalPages)
        return lastPage.currentPage + 1;
      return undefined;
    },
    // staleTime: 3000,
  });
};

export default useGetMyMarker;
