import { useEffect } from "react";
import type { Marker } from "../Map/Map";
import noimg from "../../assets/images/noimg.png";
import * as Styled from "./MarkerInfoModal.style";

interface Props {
  currentMarkerInfo: Marker;
}

const MarkerInfoModal = ({ currentMarkerInfo }: Props) => {
  useEffect(() => {
    if (currentMarkerInfo.photos) {
      console.log(currentMarkerInfo.photos[0].photoUrl);
    }
  }, []);
  return (
    <div>
      <Styled.imageWrap>
        <img
          src={
            currentMarkerInfo.photos
              ? currentMarkerInfo.photos[0].photoUrl
              : noimg
          }
          alt=""
          width={"80%"}
          height={250}
        />
      </Styled.imageWrap>
      <Styled.description>{currentMarkerInfo.description}</Styled.description>
    </div>
  );
};

export default MarkerInfoModal;
