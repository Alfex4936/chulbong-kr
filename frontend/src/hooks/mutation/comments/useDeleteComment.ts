import { useMutation, useQueryClient } from "@tanstack/react-query";
import deleteComment from "../../../api/comments/deleteComment";

const useDeleteComment = (markerId: number) => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (commentId: number) => {
      return deleteComment(commentId);
    },
    onSuccess: () => {
      queryClient.resetQueries({ queryKey: ["comments", markerId] });
    },
    onError: () => {
      alert("잠시 후 다시 시도해 주세요!");
    },
  });
};

export default useDeleteComment;
