import getAllReports from "@/api/report/getAllReports";
import { useQuery } from "@tanstack/react-query";

const useAllRerportData = () => {
  return useQuery({
    queryKey: ["marker", "report", "all"],
    queryFn: getAllReports,
  });
};

export default useAllRerportData;
