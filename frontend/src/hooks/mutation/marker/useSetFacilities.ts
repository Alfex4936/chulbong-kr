import setFacilities, { type Facilities } from "@/api/markers/setFacilities";
import { useMutation } from "@tanstack/react-query";

const useSetFacilities = () => {
  return useMutation({
    mutationFn: (props: { markerId: number; facilities: Facilities[] }) => {
      return setFacilities(props.markerId, props.facilities);
    },
  });
};

export default useSetFacilities;
