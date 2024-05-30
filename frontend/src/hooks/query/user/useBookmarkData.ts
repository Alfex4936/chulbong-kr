import bookmarkMarker from "@/api/user/bookmarkMarker";
import { useQuery } from "@tanstack/react-query";

const useBookmarkData = () => {
  return useQuery({
    queryKey: ["marker", "bookmark"],
    queryFn: bookmarkMarker,
  });
};

export default useBookmarkData;
