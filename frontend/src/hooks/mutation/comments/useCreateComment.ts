import createComment from "@/api/comments/createComment";
import { useToast } from "@/components/ui/use-toast";
import useLoginModalStateStore from "@/store/useLoginModalStateStore";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { isAxiosError } from "axios";

const useCreateComment = (
  body: {
    markerId: number;
    commentText: string;
  },
  callback?: VoidFunction
) => {
  const queryClient = useQueryClient();

  const { open: openLoginModal } = useLoginModalStateStore();
  const { toast } = useToast();

  return useMutation({
    mutationFn: () => {
      return createComment(body);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["comments", body.markerId] });
      if (callback) callback();
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

export default useCreateComment;
