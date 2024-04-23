import { useMutation } from "@tanstack/react-query";
import deleteUser from "@/api/user/deleteUser";

const useDeleteUser = () => {
  return useMutation({
    mutationFn: deleteUser,
  });
};

export default useDeleteUser;
