import search from "@/api/search/search";
import { useQuery } from "@tanstack/react-query";

const useSearch = (query: string) => {
  return useQuery({
    queryKey: ["search", "marker", query],
    queryFn: () => {
      return search(query);
    },
    enabled: query.length > 0,
  });
};

export default useSearch;
