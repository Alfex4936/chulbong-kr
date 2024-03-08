import { useInfiniteQuery } from "@tanstack/react-query";
import getCommets from "../../../api/comments/getCommets";

const useGetComments = (id: number) => {
  return useInfiniteQuery({
    queryKey: ["comments", id],
    queryFn: ({ pageParam }) => {
      return getCommets({ id, pageParam });
    },
    initialPageParam: 1,
    getNextPageParam: (lastPage) => {
      if (lastPage.currentPage < lastPage.totalPages)
        return lastPage.currentPage + 1;
      return undefined;
    },
  });
};

export default useGetComments;
