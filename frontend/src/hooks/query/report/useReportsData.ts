import getAllReports from "@/api/report/getAllReports";
import getMyReports from "@/api/report/getMyReports";
import getReportsForMarker from "@/api/report/getReportsForMarker";
import { useQuery } from "@tanstack/react-query";

interface Props {
  markerId?: number;
  type?: string;
}

const useReportsData = ({ markerId, type = "me" }: Props) => {
  if (type === "me") {
    return useQuery({
      queryKey: ["marker", "report", type],
      queryFn: getMyReports,
    });
  } else if (type === "all") {
    return useQuery({
      queryKey: ["marker", "report", type],
      queryFn: getAllReports,
    });
  } else if (type === "formarker") {
    console.log(markerId);
    return useQuery({
      queryKey: ["marker", "report", type, markerId],
      queryFn: () => {
        return getReportsForMarker(markerId as number);
      },
    });
  }
};

export default useReportsData;
