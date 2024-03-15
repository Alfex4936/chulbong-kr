import { useMutation, useQueryClient } from "@tanstack/react-query";
import logout from "../../../api/auth/logout";

const useLogout = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: logout,
    onSuccess: () => {
      queryClient.removeQueries({ queryKey: ["myInfo"] });
      queryClient.removeQueries({ queryKey: ["favorite"] });
      queryClient.removeQueries({ queryKey: ["dislikeState"] });
      queryClient.removeQueries({ queryKey: ["myMarker"] });
    },
  });
};

export default useLogout;
