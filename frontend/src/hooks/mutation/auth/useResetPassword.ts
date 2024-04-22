import resetPassword from "@/api/auth/resetPassword";
import { useMutation } from "@tanstack/react-query";
import { useRouter } from "next/navigation";
import logout from "@/api/auth/logout";

const useResetPassword = () => {
  const router = useRouter();
  return useMutation({
    mutationFn: resetPassword,
    onSuccess: async () => {
      await logout();
      router.push("/signin");
    },
  });
};

export default useResetPassword;
