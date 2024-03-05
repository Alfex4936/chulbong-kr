import { useMutation } from "@tanstack/react-query";
import requestPasswordReset from "../../../api/auth/requestPasswordReset";

const useRequestPasswordReset = (email: string) => {
  return useMutation({
    mutationFn: () => {
      return requestPasswordReset(email);
    },
  });
};

export default useRequestPasswordReset;
