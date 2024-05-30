import { useInfiniteQuery } from "@tanstack/react-query";
import mylocateMarker from "@/api/user/mylocateMarker";

const useMylocateData = () => {
  return useInfiniteQuery({
    queryKey: ["myMarker"],
    queryFn: async ({ pageParam }) => {
      return mylocateMarker({ pageParam });
    },
    initialPageParam: 1,
    getNextPageParam: (lastPage) => {
      if (lastPage.currentPage < lastPage.totalPages)
        return lastPage.currentPage + 1;
      return undefined;
    },
    retry: false,
  });
};

export default useMylocateData;
