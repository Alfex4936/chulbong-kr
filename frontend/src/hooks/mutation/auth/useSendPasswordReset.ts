import sendPasswordReset from "@/api/auth/sendPasswordReset";
import { useMutation } from "@tanstack/react-query";

const useSendPasswordReset = () => {
  return useMutation({
    mutationFn: sendPasswordReset,
  });
};

export default useSendPasswordReset;
