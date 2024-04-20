import logout from "@/api/auth/logout";
import { useMutation } from "@tanstack/react-query";

const useLogout = () => {
  return useMutation({
    mutationFn: logout,
    onError: (error) => {
      console.log(error);
    },
    onSuccess: () => {
      window.location.reload();
    },
  });
};

export default useLogout;
