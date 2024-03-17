import { useQuery } from "@tanstack/react-query";
import getFacilities from "../../../api/markers/getFacilities";

const useGetFacilities = (id: number) => {
  return useQuery({
    queryKey: ["facilities", id],
    queryFn: () => {
      return getFacilities(id);
    },
    retry: false,
  });
};

export default useGetFacilities;
