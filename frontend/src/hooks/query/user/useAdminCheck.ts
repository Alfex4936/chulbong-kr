import { useQuery } from "@tanstack/react-query";
import adminCheck from "../../../api/user/adminCheck";

const useAdminCheck = () => {
  return useQuery({
    queryKey: ["adminCheck"],
    queryFn: adminCheck,
    retry: false,
  });
};

export default useAdminCheck;
