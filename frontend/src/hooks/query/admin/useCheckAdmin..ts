import checkAdmin from "@/api/admin/checkAdmin";
import { useQuery } from "@tanstack/react-query";

const useCheckAdmin = () => {
  return useQuery({
    queryKey: ["admin", "check"],
    queryFn: checkAdmin,
  });
};

export default useCheckAdmin;
