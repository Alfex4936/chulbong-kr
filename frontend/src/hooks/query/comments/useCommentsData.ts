import { useInfiniteQuery } from "@tanstack/react-query";
import getComments from "@/api/comments/getComments";

const useGetComments = (id: number) => {
  return useInfiniteQuery({
    queryKey: ["comments", id],
    queryFn: ({ pageParam }) => {
      return getComments({ id, pageParam });
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
