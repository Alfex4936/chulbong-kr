import { useEffect, useState } from "react";
import useMapPositionStore from "../store/useMapPositionStore";
import getAddress, { type AddressInfo } from "../utils/getAddress";

interface AddressViewInfo {
  depth1: string;
  depth2: string;
  depth3: string;
}

const useAddressData = () => {
  const positionState = useMapPositionStore();
  const [address, setAddress] = useState<AddressViewInfo | null>(null);
  const [isError, setIsError] = useState(false);

  useEffect(() => {
    const fetch = async () => {
      try {
        const data = (await getAddress(
          positionState.lat,
          positionState.lng
        )) as AddressInfo;

        setAddress({
          depth1: data.region_1depth_name,
          depth2: data.region_2depth_name,
          depth3: data.region_3depth_name,
        });
      } catch {
        setIsError(true);
      }
    };

    fetch();
  }, [positionState.lat, positionState.lng]);

  return { address, isError };
};

export default useAddressData;
