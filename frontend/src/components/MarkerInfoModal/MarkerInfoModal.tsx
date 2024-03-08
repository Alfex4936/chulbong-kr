import DeleteOutlineIcon from "@mui/icons-material/DeleteOutline";
import RateReviewIcon from "@mui/icons-material/RateReview";
import ThumbDownAltIcon from "@mui/icons-material/ThumbDownAlt";
import ThumbDownOffAltIcon from "@mui/icons-material/ThumbDownOffAlt";
import CircularProgress from "@mui/material/CircularProgress";
import IconButton from "@mui/material/IconButton";
import Tooltip from "@mui/material/Tooltip";
import { isAxiosError } from "axios";
import { useEffect } from "react";
import noimg from "../../assets/images/noimg.webp";
import useDeleteMarker from "../../hooks/mutation/marker/useDeleteMarker";
import useMarkerDislike from "../../hooks/mutation/marker/useMarkerDislike";
import useUndoDislike from "../../hooks/mutation/marker/useUndoDislike";
import useDislikeState from "../../hooks/query/marker/useDislikeState";
import useGetMarker from "../../hooks/query/marker/useGetMarker";
import useModalStore from "../../store/useModalStore";
import useToastStore from "../../store/useToastStore";
import useUserStore from "../../store/useUserStore";
import type { MarkerClusterer } from "../../types/Cluster.types";
import type { KakaoMarker } from "../../types/KakaoMap.types";
import type { MarkerInfo } from "../Map/Map";
import * as Styled from "./MarkerInfoModal.style";
import MarkerInfoSkeleton from "./MarkerInfoSkeleton";

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

  const {
    data: dislikeState,
    isError: isDislikeError,
    isLoading: dislikeLoading,
  } = useDislikeState(currentMarkerInfo.markerId);

  const { mutateAsync: doDislike, isPending: disLikePending } =
    useMarkerDislike(currentMarkerInfo.markerId);

  const { mutateAsync: undoDislike, isPending: undoDislikePending } =
    useUndoDislike(currentMarkerInfo.markerId);

  const { mutateAsync: deleteMarker, isPending: deleteLoading } =
    useDeleteMarker(currentMarkerInfo.markerId);

  useEffect(() => {
    toastState.close();
    toastState.setToastText("");
  }, []);

  const handleDelete = async () => {
    try {
      await deleteMarker();

      const newMarkers = markers.filter(
        (_, index) => index !== currentMarkerInfo.index
      );

      markers[currentMarkerInfo.index].setMap(null);
      setMarkers(newMarkers);

      clusterer.removeMarker(markers[currentMarkerInfo.index]);
      setMarkerInfoModal(false);

      toastState.setToastText("삭제 완료");
      toastState.open();
    } catch (error) {
      alert("삭제 실패 잠시 후 다시 시도해주세요!");
    }
  };

  const handleViewReview = () => {
    console.log("리뷰 보기");
  };

  const handleDislike = async () => {
    if (disLikePending || undoDislikePending) return;
    try {
      if (dislikeState?.disliked) {
        await undoDislike();
      } else {
        await doDislike();
      }
    } catch (error) {
      if (isAxiosError(error)) {
        if (error.response?.status === 401) {
          setMarkerInfoModal(false);
          modalState.openLogin();
        }
      }
    }
  };

  if (isLoading) return <MarkerInfoSkeleton />;
  if (isError) return <div>에러</div>;

  return (
    <div>
      <Styled.imageWrap>
        <img
          src={marker?.photos ? marker.photos[0].photoUrl : noimg}
          alt="철봉 상세 이미지"
        />
        <Styled.description>{marker?.description}</Styled.description>
      </Styled.imageWrap>
      <Styled.BottomButtons>
        <Tooltip title="리뷰 보기" arrow disableInteractive>
          <IconButton onClick={handleViewReview} aria-label="review">
            <RateReviewIcon />
          </IconButton>
        </Tooltip>
        {dislikeState?.disliked && !isDislikeError ? (
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
                  <Styled.DislikeCount>
                    {marker?.dislikeCount || 0}
                  </Styled.DislikeCount>
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
              )}
            </IconButton>
          </Tooltip>
        )}

        {userState.user.user.userId === marker?.userId && (
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
