import { useQuery } from "@tanstack/react-query";
import getDislikeState from "../../../api/markers/getDislikeState";

const useDislikeState = (id: number) => {
  return useQuery({
    queryKey: ["dislikeState", id],
    queryFn: () => {
      return getDislikeState(id);
    },
    retry: false,
  });
};

export default useDislikeState;
