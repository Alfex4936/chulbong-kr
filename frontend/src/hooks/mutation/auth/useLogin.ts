import login from "@/api/auth/login";
import { useMutation } from "@tanstack/react-query";

const useLogin = () => {
  return useMutation({
    mutationFn: login,
    onError: (error) => {
      console.log(error);
    },
  });
};

export default useLogin;
