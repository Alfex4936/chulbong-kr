import getFacilities from "@/api/markers/getFacilities";
import { useQuery } from "@tanstack/react-query";

const useFacilitiesData = (id: number) => {
  return useQuery({
    queryKey: ["facilities", id],
    queryFn: () => {
      return getFacilities(id);
    },
  });
};

export default useFacilitiesData;
