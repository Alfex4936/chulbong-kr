import { useMutation, useQueryClient } from "@tanstack/react-query";
import { isAxiosError } from "axios";
import createComment from "../../../api/comments/createComment";
import useModalStore from "../../../store/useModalStore";

const useCreateComment = (body: { markerId: number; commentText: string }) => {
  const modalState = useModalStore();
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: () => {
      return createComment(body);
    },
    onSuccess: () => {
      queryClient.resetQueries({ queryKey: ["comments", body.markerId] });
    },
    onError: (error) => {
      if (isAxiosError(error)) {
        if (error?.response?.status === 401) {
          modalState.openLogin();
        }
      }
    },
  });
};

export default useCreateComment;
