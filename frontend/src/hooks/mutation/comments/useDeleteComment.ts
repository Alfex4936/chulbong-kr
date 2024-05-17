import deleteComment from "@/api/comments/deleteComment";
import { useToast } from "@/components/ui/use-toast";
import useLoginModalStateStore from "@/store/useLoginModalStateStore";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { isAxiosError } from "axios";

const useDeleteComment = (markerId: number) => {
  const queryClient = useQueryClient();

  
  const { open: openLoginModal } = useLoginModalStateStore();
  const { toast } = useToast();

  return useMutation({
    mutationFn: (commentId: number) => {
      return deleteComment(commentId);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["comments", markerId] });
    },
    onError: (error) => {
      if (isAxiosError(error)) {
        if (error?.response?.status === 401) {
          openLoginModal();
        } else {
          toast({ description: "잠시 후 다시 시도해주세요." });
        }
      } else {
        toast({ description: "잠시 후 다시 시도해주세요." });
      }
    },
  });
};

export default useDeleteComment;
