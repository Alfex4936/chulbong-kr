import updateUserName from "@/api/user/updateUserName";
import { useToast } from "@/components/ui/use-toast";
import { useMutation, useQueryClient } from "@tanstack/react-query";

const useUpdateUserName = () => {
  const queryClient = useQueryClient();
  const { toast } = useToast();

  return useMutation({
    mutationFn: updateUserName,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["user", "me"] });
    },
    onError: () => {
      toast({ description: "잠시 후 다시 시도해 주세요." });
    },
  });
};

export default useUpdateUserName;
