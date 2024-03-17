import { useMutation } from "@tanstack/react-query";
import setFacilities, {
  type Facilities,
} from "../../../api/markers/setFacilities";

const useSetFacilities = () => {
  return useMutation({
    mutationFn: (props: { markerId: number; facilities: Facilities[] }) => {
      return setFacilities(props.markerId, props.facilities);
    },
  });
};

export default useSetFacilities;
