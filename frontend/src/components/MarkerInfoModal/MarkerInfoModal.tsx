import RemoveIcon from "@mui/icons-material/Remove";
import CircularProgress from "@mui/material/CircularProgress";
import IconButton from "@mui/material/IconButton";
import Tooltip from "@mui/material/Tooltip";
import { useEffect, useState } from "react";
import deleteMarker from "../../api/markers/deleteMarker";
import noimg from "../../assets/images/noimg.webp";
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

  const [loading, setLoading] = useState(false);

  useEffect(() => {
    toastState.close();
    toastState.setToastText("");
  }, []);

  const handleDelete = () => {
    setLoading(true);
    deleteMarker(currentMarkerInfo.markerId)
      .then(() => {
        toastState.setToastText("삭제 완료");
        toastState.open();

        markers[currentMarkerInfo.index].setMap(null);
        setMarkerInfoModal(false);
      })
      .catch((error) => {
        // 임시(확인용)
        console.log(error);
        alert("삭제 실패 잠시 후 다시 시도해주세요!");
      })
      .finally(() => {
        setLoading(false);
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
            {loading ? (
              <CircularProgress color="inherit" size={20} />
            ) : (
              <RemoveIcon />
            )}
          </IconButton>
        </Tooltip>
      )}

      <Styled.imageWrap>
        <img
          src={currentMarkerInfo.photos ? currentMarkerInfo.photos[0] : noimg}
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
