import updateUserName from "@/api/user/updateUserName";
import { useMutation, useQueryClient } from "@tanstack/react-query";

const useUpdateUserName = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: updateUserName,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["user", "me"] });
    },
  });
};

export default useUpdateUserName;
