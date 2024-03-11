// import { useMutation, useQueryClient } from "@tanstack/react-query";
// import signin from "@/api/auth/signin";


// const useLogout = () => {
//   const queryClient = useQueryClient();

//   return useMutation({
//     mutationFn: logout,
//     onSuccess: () => {
//       queryClient.removeQueries({ queryKey: ["myInfo"] });
//       queryClient.removeQueries({ queryKey: ["dislikeState"] });
//     },
//   });
// };

// export default useLogout;
