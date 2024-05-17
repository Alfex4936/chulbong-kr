import useMapStatusStore from "@/store/useMapStatusStore";
import useMapStore from "@/store/useMapStore";
import getAddress, { type AddressInfo } from "@/utils/getAddress";
import { useEffect, useState } from "react";

interface AddressViewInfo {
  depth1: string;
  depth2: string;
  depth3: string;
}

const useAddressData = () => {
  const positionState = useMapStatusStore();
  const [address, setAddress] = useState<AddressViewInfo | null>(null);
  const [isError, setIsError] = useState(false);
  const { map } = useMapStore();

  useEffect(() => {
    if (!map) return;
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
      } catch (error) {
        console.log(error);
        setIsError(true);
      }
    };

    fetch();
  }, [positionState.lat, positionState.lng, map]);

  return { address, isError };
};

export default useAddressData;
