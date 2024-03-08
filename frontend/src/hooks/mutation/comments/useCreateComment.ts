import { useMutation, useQueryClient } from "@tanstack/react-query";
import createComment from "../../../api/comments/createComment";

const useCreateComment = (body: { markerId: number; commentText: string }) => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: () => {
      return createComment(body);
    },
    onSuccess: () => {
      queryClient.resetQueries({ queryKey: ["comments", body.markerId] });
    },
    onError: () => {
      alert("잠시 후 다시 시도해 주세요!");
    },
  });
};

export default useCreateComment;
