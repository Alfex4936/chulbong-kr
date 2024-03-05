import { useMutation } from "@tanstack/react-query";
import resetPassword from "../../../api/auth/resetPassword";

const useResetPassword = (token: string, password: string) => {
  return useMutation({
    mutationFn: () => {
      return resetPassword(token, password);
    },
    onSuccess(data) {
      console.log(data);
    },
    onError(err) {
      console.log(err);
    },
  });
};

export default useResetPassword;
