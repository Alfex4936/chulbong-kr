import { useMutation, useQueryClient } from "@tanstack/react-query";
import deleteUser from "@/api/user/deleteUser";
import { useRouter } from "next/navigation";
import { useToast } from "@/components/ui/use-toast";

const useDeleteUser = () => {
  const router = useRouter();
  const queryClient = useQueryClient();
  const { toast } = useToast();

  return useMutation({
    mutationFn: deleteUser,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["user", "me"] });
      router.push("/home");
    },
    onError: () => {
      toast({ description: "잠시 후 다시 시도해주세요." });
    },
  });
};

export default useDeleteUser;
