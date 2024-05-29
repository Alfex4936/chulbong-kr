import logout from "@/api/auth/logout";
import { useToast } from "@/components/ui/use-toast";
import { useMutation, useQueryClient } from "@tanstack/react-query";

const useLogout = () => {
  const { toast } = useToast();

  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: logout,
    onError: () => {
      toast({ description: "잠시 후 다시 시도해주세요." });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["user", "me"] });
      window.location.reload();
    },
  });
};

export default useLogout;
