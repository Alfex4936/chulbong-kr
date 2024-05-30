import myInfo from "@/api/user/myInfo";
import { useQuery } from "@tanstack/react-query";

const useMyinfoData = () => {
  return useQuery({
    queryKey: ["user", "me"],
    queryFn: myInfo,
  });
};

export default useMyinfoData;
