import { useQuery } from "@tanstack/react-query";
import getMyInfo from "../../../api/user/getMyInfo";

const useGetMyInfo = () => {
  return useQuery({
    queryKey: ["myInfo"],
    queryFn: getMyInfo,
    retry: false,
  });
};

export default useGetMyInfo;
