import getSearchLoation from "@/api/kakao/getSearchLocation";
import { useQuery } from "@tanstack/react-query";

const useSearchLocationData = (query: string) => {
  return useQuery({
    queryKey: ["search", query],
    queryFn: () => {
      return getSearchLoation(query);
    },
    enabled: query.length > 0,
  });
};

export default useSearchLocationData;
