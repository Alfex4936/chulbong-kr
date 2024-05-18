import sendVerifyCode from "@/api/auth/sendVerifyCode";
import { useMutation } from "@tanstack/react-query";

const useSendVerifyCode = () => {
  return useMutation({
    mutationFn: sendVerifyCode,
  });
};

export default useSendVerifyCode;
