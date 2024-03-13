import { useMutation } from "@tanstack/react-query";
import updateDescription from "../../../api/markers/updateDescription";
import { useQueryClient } from "@tanstack/react-query";

const useUpdateDesc = (desc: string, id: number) => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: () => {
      return updateDescription(desc, id);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["marker", id] });
    },
  });
};

export default useUpdateDesc;
