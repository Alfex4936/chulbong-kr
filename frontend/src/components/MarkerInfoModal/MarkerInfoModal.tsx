import DeleteOutlineIcon from "@mui/icons-material/DeleteOutline";
import RateReviewIcon from "@mui/icons-material/RateReview";
import CircularProgress from "@mui/material/CircularProgress";
import IconButton from "@mui/material/IconButton";
import Tooltip from "@mui/material/Tooltip";
import { useEffect, useState } from "react";
import deleteMarker from "../../api/markers/deleteMarker";
import noimg from "../../assets/images/noimg.webp";
import useToastStore from "../../store/useToastStore";
import useUserStore from "../../store/useUserStore";
import type { KakaoMarker } from "../../types/KakaoMap.types";
import type { MarkerInfo } from "../Map/Map";
import * as Styled from "./MarkerInfoModal.style";
import ThumbDownAltIcon from "@mui/icons-material/ThumbDownAlt";
import type { MarkerClusterer } from "../../types/Cluster.types";

interface Props {
  currentMarkerInfo: MarkerInfo;
  setMarkerInfoModal: React.Dispatch<React.SetStateAction<boolean>>;
  markers: KakaoMarker[];
  setMarkers: React.Dispatch<React.SetStateAction<KakaoMarker[]>>;
  clusterer: MarkerClusterer;
}

const MarkerInfoModal = ({
  currentMarkerInfo,
  setMarkerInfoModal,
  markers,
  setMarkers,
  clusterer,
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

        const newMarkers = markers.filter(
          (_, index) => index !== currentMarkerInfo.index
        );

        markers[currentMarkerInfo.index].setMap(null);
        setMarkers(newMarkers);

        clusterer.removeMarker(markers[currentMarkerInfo.index]);

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

  const handleViewReview = () => {
    console.log("리뷰 보기");
  };

  const handleDislike = () => {
    console.log("싫어요");
  };

  return (
    <div>
      <Styled.imageWrap>
        <img
          src={currentMarkerInfo.photos ? currentMarkerInfo.photos[0] : noimg}
          alt=""
          width={"90%"}
          height={300}
        />
        <Styled.description>{currentMarkerInfo.description}</Styled.description>
      </Styled.imageWrap>
      <Styled.BottomButtons>
        <Tooltip title="리뷰 보기" arrow disableInteractive>
          <IconButton onClick={handleViewReview} aria-label="review">
            <RateReviewIcon />
          </IconButton>
        </Tooltip>
        <Tooltip title="싫어요" arrow disableInteractive>
          <IconButton onClick={handleDislike} aria-label="dislike">
            <ThumbDownAltIcon />
          </IconButton>
        </Tooltip>
        {userState.user.user.userId === currentMarkerInfo.userId && (
          <Tooltip title="삭제 하기" arrow disableInteractive>
            <IconButton onClick={handleDelete} aria-label="delete">
              {loading ? (
                <CircularProgress color="inherit" size={20} />
              ) : (
                <DeleteOutlineIcon />
              )}
            </IconButton>
          </Tooltip>
        )}
      </Styled.BottomButtons>
    </div>
  );
};

export default MarkerInfoModal;
