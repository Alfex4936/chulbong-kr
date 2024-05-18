import login from "@/api/auth/login";
import { useMutation } from "@tanstack/react-query";
import { useRouter, useSearchParams } from "next/navigation";

const useLogin = () => {
  const router = useRouter();
  const searchParams = useSearchParams();
  const redirect = searchParams.get("redirect");

  return useMutation({
    mutationFn: login,
    onSuccess: () => {
      router.push(redirect || "/home");
    },
  });
};

export default useLogin;
