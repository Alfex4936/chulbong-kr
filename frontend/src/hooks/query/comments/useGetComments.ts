import { useQuery } from "@tanstack/react-query";
import getCommets from "../../../api/comments/getCommets";

const useGetComments = (id: number) => {
  return useQuery({
    queryKey: ["comments", id],
    queryFn: () => {
      return getCommets(id);
    },
    retry: false,
  });
};

export default useGetComments;
