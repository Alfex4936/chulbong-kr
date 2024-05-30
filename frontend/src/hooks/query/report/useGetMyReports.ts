import getMyReports from "@/api/report/getMyReports";
import { useQuery } from "@tanstack/react-query";

const useGetMyReports = () => {
  return useQuery({
    queryKey: ["marker", "report"],
    queryFn: getMyReports,
  });
};

export default useGetMyReports;
