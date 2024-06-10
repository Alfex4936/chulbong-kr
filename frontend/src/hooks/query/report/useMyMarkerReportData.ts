import getReportForMyMarker from "@/api/report/getReportForMyMarker";
import { useQuery } from "@tanstack/react-query";

const useMyMarkerReportData = () => {
  return useQuery({
    queryKey: ["report", "my"],
    queryFn: getReportForMyMarker,
  });
};

export default useMyMarkerReportData;
