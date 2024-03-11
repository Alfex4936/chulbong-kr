import { useMutation, useQueryClient } from "@tanstack/react-query";
import login from "../../../api/auth/login";

const useSignin = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: login,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["myInfo"] });
      queryClient.invalidateQueries({ queryKey: ["dislikeState"] });
      queryClient.invalidateQueries({ queryKey: ["favorite"] });
    },
  });
};

export default useSignin;
