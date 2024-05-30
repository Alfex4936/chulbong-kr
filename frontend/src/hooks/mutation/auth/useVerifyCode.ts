import verifyCode from "@/api/auth/verifyCode";
import { useMutation } from "@tanstack/react-query";

const useVerifyCode = () => {
  return useMutation({
    mutationFn: verifyCode,
  });
};

export default useVerifyCode;
