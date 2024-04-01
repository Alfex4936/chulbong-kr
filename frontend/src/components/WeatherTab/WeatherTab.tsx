import ArrowForwardIosIcon from "@mui/icons-material/ArrowForwardIos";
import { useEffect, useState } from "react";
import useMapPositionStore from "../../store/useMapPositionStore";
import getAddress, { type AddressInfo } from "../../utils/getAddress";
import * as Styled from "./WeatherTab.style";

interface AddressViewInfo {
  depth1: string;
  depth2: string;
  depth3: string;
}

const WeatherTab = () => {
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

  return (
    <Styled.Container>
      {!isError && address?.depth1 ? (
        <Styled.WeatherWrap>
          <div>{address.depth1}</div>
          <span>
            <ArrowForwardIosIcon sx={{ fontSize: "1rem", color: "#888" }} />
          </span>
          <div>{address.depth2}</div>
          <span>
            <ArrowForwardIosIcon sx={{ fontSize: "1rem", color: "#888" }} />
          </span>
          <div>{address.depth3}</div>
        </Styled.WeatherWrap>
      ) : (
        <div>주소가 지원되지 않습니다.</div>
      )}
    </Styled.Container>
  );
};

export default WeatherTab;
