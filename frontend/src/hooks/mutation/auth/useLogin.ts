import login from "@/api/auth/login";
import { useMutation } from "@tanstack/react-query";
import { useRouter } from "next/navigation";

const useLogin = () => {
  const router = useRouter();
  return useMutation({
    mutationFn: login,
    onError: (error) => {
      console.log(error);
    },
    onSuccess: () => {
      router.push("/home");
    },
  });
};

export default useLogin;
