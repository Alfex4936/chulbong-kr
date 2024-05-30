import getAllReports from "@/api/report/getAllReports";
import getMyReports from "@/api/report/getMyReports";
import getReportsForMarker from "@/api/report/getReportsForMarker";
import { useQuery } from "@tanstack/react-query";

interface Props {
  markerId?: number;
  type?: string;
}

const useReportsData = ({ markerId, type = "me" }: Props) => {
  let queryFn;
  let queryKey: (string | number)[] = ["marker", "report", type];

  switch (type) {
    case "me":
      queryFn = getMyReports;
      break;
    case "all":
      queryFn = getAllReports;
      break;
    case "formarker":
      if (markerId === undefined) {
        throw new Error(`markerId 없음`);
      }
      queryFn = () => getReportsForMarker(markerId);
      queryKey.push(markerId);
      break;
    default:
      throw new Error(`존재하지 않는 타입: ${type}`);
  }

  return useQuery({
    queryKey,
    queryFn,
    select: (data) => {
      const sortedData = [...data].sort(
        (a, b) =>
          new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime()
      );
      return sortedData;
    },
  });
};

export default useReportsData;
