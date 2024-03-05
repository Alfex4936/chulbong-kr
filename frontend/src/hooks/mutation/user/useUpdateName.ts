import { useMutation, useQueryClient } from "@tanstack/react-query";
import updateName from "../../../api/user/updateName";

const useUpdateName = (name: string) => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: () => {
      return updateName(name);
    },
    onSuccess(data) {
      queryClient.invalidateQueries({ queryKey: ["myInfo"] });
      console.log(data);
    },
    onError(err) {
      console.log(err);
    },
  });
};

export default useUpdateName;
