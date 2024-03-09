import CloseIcon from "@mui/icons-material/Close";
import DeleteOutlineIcon from "@mui/icons-material/DeleteOutline";
import RateReviewIcon from "@mui/icons-material/RateReview";
import ThumbDownAltIcon from "@mui/icons-material/ThumbDownAlt";
import ThumbDownOffAltIcon from "@mui/icons-material/ThumbDownOffAlt";
import CircularProgress from "@mui/material/CircularProgress";
import IconButton from "@mui/material/IconButton";
import Tooltip from "@mui/material/Tooltip";
import { isAxiosError } from "axios";
import { useEffect, useState } from "react";
import noimg from "../../assets/images/noimg.webp";
import useDeleteMarker from "../../hooks/mutation/marker/useDeleteMarker";
import useMarkerDislike from "../../hooks/mutation/marker/useMarkerDislike";
import useUndoDislike from "../../hooks/mutation/marker/useUndoDislike";
import useGetMarker from "../../hooks/query/marker/useGetMarker";
import useModalStore from "../../store/useModalStore";
import useToastStore from "../../store/useToastStore";
import useUserStore from "../../store/useUserStore";
import type { MarkerClusterer } from "../../types/Cluster.types";
import type { KakaoMarker } from "../../types/KakaoMap.types";
import type { MarkerInfo } from "../Map/Map";
import * as Styled from "./MarkerInfoModal.style";
import MarkerInfoSkeleton from "./MarkerInfoSkeleton";
import MarkerReview from "./MarkerReview";

interface Props {
  currentMarkerInfo: MarkerInfo;
  setMarkerInfoModal: React.Dispatch<React.SetStateAction<boolean>>;
  markers: KakaoMarker[];
  setMarkers: React.Dispatch<React.SetStateAction<KakaoMarker[]>>;
  clusterer: MarkerClusterer;
}

const MarkerInfoModal = ({
  markers,
  setMarkers,
  clusterer,
  currentMarkerInfo,
  setMarkerInfoModal,
}: Props) => {
  const userState = useUserStore();
  const toastState = useToastStore();
  const modalState = useModalStore();

  const {
    data: marker,
    isLoading,
    isError,
  } = useGetMarker(currentMarkerInfo.markerId);

  const { mutateAsync: doDislike, isPending: disLikePending } =
    useMarkerDislike(currentMarkerInfo.markerId);

  const { mutateAsync: undoDislike, isPending: undoDislikePending } =
    useUndoDislike(currentMarkerInfo.markerId);

  const { mutateAsync: deleteMarker } = useDeleteMarker(
    currentMarkerInfo.markerId
  );

  const [isReview, setIsReview] = useState(false);

  const [deleteLoading, setDeleteLoading] = useState(false);

  useEffect(() => {
    toastState.close();
    toastState.setToastText("");
  }, []);

  const filtering = async () => {
    const marker = markers.find(
      (value) => Number(value.Gb) === currentMarkerInfo.markerId
    );

    const newMarkers = markers.filter(
      (value) => Number(value.Gb) !== currentMarkerInfo.markerId
    );

    if (marker) {
      marker.setMap(null);
      clusterer.removeMarker(marker);
      setMarkers(newMarkers);
    }
  };

  const handleDelete = async () => {
    setDeleteLoading(true);
    try {
      await deleteMarker();
      await filtering();
      setMarkerInfoModal(false);
      toastState.setToastText("삭제 완료");
      toastState.open();
    } catch (error) {
      alert("삭제 실패 잠시 후 다시 시도해주세요!");
    } finally {
      setDeleteLoading(false);
    }
  };

  const handleViewReview = () => {
    setIsReview(true);
  };

  const handleDislike = async () => {
    if (disLikePending || undoDislikePending) return;
    try {
      if (marker?.disliked) {
        await undoDislike();
      } else {
        await doDislike();
      }
    } catch (error) {
      if (isAxiosError(error)) {
        if (error.response?.status === 401) {
          modalState.openLogin();
        }
      }
    }
  };

  if (isLoading) return <MarkerInfoSkeleton />;
  if (isError)
    return (
      <div style={{ fontSize: "1.2rem" }}>
        <Tooltip title="닫기" arrow disableInteractive>
          <IconButton
            onClick={() => {
              setMarkerInfoModal(false);
            }}
            aria-label="delete"
            sx={{
              position: "absolute",
              top: "0",
              right: "0",
            }}
          >
            <CloseIcon />
          </IconButton>
        </Tooltip>
        잠시 후 다시 시도해 주세요.....
      </div>
    );

  return (
    <div>
      {isReview ? (
        <MarkerReview
          setIsReview={setIsReview}
          markerId={marker?.markerId as number}
        />
      ) : (
        <>
          <Tooltip title="닫기" arrow disableInteractive>
            <IconButton
              onClick={() => {
                setMarkerInfoModal(false);
              }}
              aria-label="delete"
              sx={{
                position: "absolute",
                top: "0",
                right: "0",
              }}
            >
              <CloseIcon />
            </IconButton>
          </Tooltip>
          <Styled.imageWrap>
            <img
              src={marker?.photos ? marker.photos[0].photoUrl : noimg}
              alt="철봉 상세 이미지"
            />
            <Styled.description>
              {marker?.description || "작성된 설명이 없습니다."}
            </Styled.description>
          </Styled.imageWrap>
          <Styled.BottomButtons>
            <Tooltip title="리뷰 보기" arrow disableInteractive>
              <IconButton onClick={handleViewReview} aria-label="review">
                <RateReviewIcon />
              </IconButton>
            </Tooltip>
            {marker?.disliked ? (
              <Tooltip title="싫어요 취소" arrow disableInteractive>
                <IconButton onClick={handleDislike} aria-label="dislike">
                  <div
                    style={{
                      width: "24px",
                      height: "24px",
                      position: "relative",
                    }}
                  >
                    <Styled.DislikeCount>
                      {marker?.dislikeCount || 0}
                    </Styled.DislikeCount>
                    <ThumbDownAltIcon />
                  </div>
                </IconButton>
              </Tooltip>
            ) : (
              <Tooltip title="싫어요" arrow disableInteractive>
                <IconButton onClick={handleDislike} aria-label="dislike">
                  <div
                    style={{
                      width: "24px",
                      height: "24px",
                      position: "relative",
                    }}
                  >
                    <Styled.DislikeCount>
                      {marker?.dislikeCount || 0}
                    </Styled.DislikeCount>
                    <ThumbDownOffAltIcon />
                  </div>
                </IconButton>
              </Tooltip>
            )}

            {(marker?.isChulbong ||
              userState.user.user.userId === marker?.userId) && (
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
        </>
      )}
    </div>
  );
};

export default MarkerInfoModal;
