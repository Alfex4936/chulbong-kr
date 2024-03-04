import DeleteOutlineIcon from "@mui/icons-material/DeleteOutline";
import RateReviewIcon from "@mui/icons-material/RateReview";
import ThumbDownAltIcon from "@mui/icons-material/ThumbDownAlt";
import ThumbDownOffAltIcon from "@mui/icons-material/ThumbDownOffAlt";
import CircularProgress from "@mui/material/CircularProgress";
import IconButton from "@mui/material/IconButton";
import Tooltip from "@mui/material/Tooltip";
import { useEffect, useState } from "react";
import deleteMarker from "../../api/markers/deleteMarker";
import getDislikeState from "../../api/markers/getDislikeState";
import markerDislike from "../../api/markers/markerDislike";
import markerUnDislike from "../../api/markers/markerUnDislike";
import noimg from "../../assets/images/noimg.webp";
import useToastStore from "../../store/useToastStore";
import useUserStore from "../../store/useUserStore";
import type { MarkerClusterer } from "../../types/Cluster.types";
import type { KakaoMarker } from "../../types/KakaoMap.types";
import type { MarkerInfo } from "../Map/Map";
import * as Styled from "./MarkerInfoModal.style";

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

  const [deleteLoading, setDeleteLoading] = useState(false);
  const [dislikeLoading, setDislikeLoading] = useState(true);

  const [disLike, setDislike] = useState(false);
  const [dislikeCount, setDislikeCount] = useState(0);

  console.log(currentMarkerInfo);

  useEffect(() => {
    const getDislike = async () => {
      try {
        const result = await getDislikeState(currentMarkerInfo.markerId);
        console.log(result);
        setDislike(result.disliked);
      } catch (error) {
        setDislike(false);
        console.log(error);
      } finally {
        setDislikeLoading(false);
      }
    };

    getDislike();
    toastState.close();
    toastState.setToastText("");
  }, []);

  useEffect(() => {
    if (currentMarkerInfo.dislikeCount) {
      setDislikeCount(currentMarkerInfo.dislikeCount);
    }
  }, []);

  const handleDelete = () => {
    setDeleteLoading(true);
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
        setDeleteLoading(false);
      });
  };

  const handleViewReview = () => {
    console.log("리뷰 보기");
  };

  const handleDislike = async () => {
    setDislikeLoading(true);
    try {
      if (disLike) {
        const result = await markerUnDislike(currentMarkerInfo.markerId);
        console.log(result);
        setDislike(false);
        setDislikeCount((prev) => prev - 1);
      } else {
        const result = await markerDislike(currentMarkerInfo.markerId);
        console.log(result);
        setDislike(true);
        setDislikeCount((prev) => prev + 1);
      }
    } catch (error) {
      console.log(error);
    } finally {
      setDislikeLoading(false);
    }
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
        {disLike ? (
          <Tooltip title="싫어요 취소" arrow disableInteractive>
            <IconButton onClick={handleDislike} aria-label="dislike">
              {dislikeLoading ? (
                <CircularProgress color="inherit" size={20} />
              ) : (
                <div
                  style={{
                    width: "24px",
                    height: "24px",
                    position: "relative",
                  }}
                >
                  <Styled.DislikeCount>{dislikeCount}</Styled.DislikeCount>
                  <ThumbDownAltIcon />
                </div>
              )}
            </IconButton>
          </Tooltip>
        ) : (
          <Tooltip title="싫어요" arrow disableInteractive>
            <IconButton onClick={handleDislike} aria-label="dislike">
              {dislikeLoading ? (
                <CircularProgress color="inherit" size={20} />
              ) : (
                <ThumbDownOffAltIcon />
              )}
            </IconButton>
          </Tooltip>
        )}

        {userState.user.user.userId === currentMarkerInfo.userId && (
          <Tooltip title="삭제 하기" arrow disableInteractive>
            <IconButton onClick={handleDelete} aria-label="delete">
              {deleteLoading ? (
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
