import RemoveIcon from "@mui/icons-material/Remove";
import { IconButton, Tooltip } from "@mui/material";
import { useEffect } from "react";
import DeleteMarker from "../../api/markers/DeleteMarker";
import noimg from "../../assets/images/noimg.png";
import useToastStore from "../../store/useToastStore";
import useUserStore from "../../store/useUserStore";
import type { KakaoMarker } from "../../types/KakaoMap.types";
import ActionButton from "../ActionButton/ActionButton";
import type { MarkerInfo } from "../Map/Map";
import * as Styled from "./MarkerInfoModal.style";

interface Props {
  currentMarkerInfo: MarkerInfo;
  setMarkerInfoModal: React.Dispatch<React.SetStateAction<boolean>>;
  markers: KakaoMarker[];
}

const MarkerInfoModal = ({
  currentMarkerInfo,
  setMarkerInfoModal,
  markers,
}: Props) => {
  const userState = useUserStore();
  const toastState = useToastStore();

  useEffect(() => {
    toastState.close();
    toastState.setToastText("");

    if (currentMarkerInfo.photos) {
      console.log(currentMarkerInfo.photos[0].photoUrl);
    }

    console.log(currentMarkerInfo);
  }, []);

  const handleDelete = () => {
    DeleteMarker(currentMarkerInfo.markerId).then((res) => {
      console.log(res);
      toastState.setToastText("삭제 완료");
      toastState.open();
      markers[currentMarkerInfo.index].setMap(null);
      setMarkerInfoModal(false);
    });
  };

  return (
    <div>
      {userState.user.user.userId === currentMarkerInfo.userId && (
        <Tooltip title="삭제하기" arrow disableInteractive>
          <IconButton
            onClick={handleDelete}
            aria-label="delete"
            sx={{
              position: "absolute",
              top: ".4rem",
              right: ".4rem",
            }}
          >
            <RemoveIcon />
          </IconButton>
        </Tooltip>
      )}

      <Styled.imageWrap>
        <img
          src={
            currentMarkerInfo.photos
              ? currentMarkerInfo.photos[0].photoUrl
              : noimg
          }
          alt=""
          width={"90%"}
          height={300}
        />
      </Styled.imageWrap>
      <Styled.description>{currentMarkerInfo.description}</Styled.description>
      <ActionButton bg="black">리뷰 보기</ActionButton>
    </div>
  );
};

export default MarkerInfoModal;
