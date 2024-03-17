import CloseIcon from "@mui/icons-material/Close";
import DeleteOutlineIcon from "@mui/icons-material/DeleteOutline";
import EditIcon from "@mui/icons-material/Edit";
import FavoriteIcon from "@mui/icons-material/Favorite";
import FavoriteBorderIcon from "@mui/icons-material/FavoriteBorder";
import PlaceIcon from "@mui/icons-material/Place";
import RateReviewIcon from "@mui/icons-material/RateReview";
import ShareIcon from "@mui/icons-material/Share";
import ThumbDownAltIcon from "@mui/icons-material/ThumbDownAlt";
import ThumbDownOffAltIcon from "@mui/icons-material/ThumbDownOffAlt";
import CircularProgress from "@mui/material/CircularProgress";
import IconButton from "@mui/material/IconButton";
import Tooltip from "@mui/material/Tooltip";
import { isAxiosError } from "axios";
import { useEffect, useState } from "react";
import { Helmet } from "react-helmet-async";
import { useNavigate } from "react-router-dom";
import noimg from "../../assets/images/noimg.webp";
import useDeleteFavorite from "../../hooks/mutation/favorites/useDeleteFavorite";
import useSetFavorite from "../../hooks/mutation/favorites/useSetFavorite";
import useDeleteMarker from "../../hooks/mutation/marker/useDeleteMarker";
import useMarkerDislike from "../../hooks/mutation/marker/useMarkerDislike";
import useUndoDislike from "../../hooks/mutation/marker/useUndoDislike";
import useUpdateDesc from "../../hooks/mutation/marker/useUpdateDesc";
import useGetMarker from "../../hooks/query/marker/useGetMarker";
import useInput from "../../hooks/useInput";
import useModalStore from "../../store/useModalStore";
import useToastStore from "../../store/useToastStore";
import type { MarkerClusterer } from "../../types/Cluster.types";
import type { KakaoMarker } from "../../types/KakaoMap.types";
import ActionButton from "../ActionButton/ActionButton";
import type { MarkerInfo } from "../Map/Map";
import * as Styled from "./MarkerInfoModal.style";
import MarkerInfoSkeleton from "./MarkerInfoSkeleton";
import MarkerReview from "./MarkerReview";
import MarkerRoadView from "./MarkerRoadView";

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
  const toastState = useToastStore();
  const modalState = useModalStore();

  const navigate = useNavigate();
  const query = new URLSearchParams(location.search);
  const sharedMarker = query.get("d");
  const sharedMarkerLat = query.get("la");
  const sharedMarkerLng = query.get("lo");

  const descInput = useInput("");

  const {
    data: marker,
    isLoading,
    isError,
    error,
  } = useGetMarker(currentMarkerInfo.markerId);

  const { mutateAsync: like, isPending: likePending } = useSetFavorite(
    currentMarkerInfo.markerId
  );
  const { mutateAsync: updateDesc } = useUpdateDesc(
    descInput.value,
    currentMarkerInfo.markerId
  );

  const { mutateAsync: unLike, isPending: unLikePending } = useDeleteFavorite(
    currentMarkerInfo.markerId
  );

  const { mutateAsync: doDislike, isPending: disLikePending } =
    useMarkerDislike(currentMarkerInfo.markerId);

  const { mutateAsync: undoDislike, isPending: undoDislikePending } =
    useUndoDislike(currentMarkerInfo.markerId);

  const { mutateAsync: deleteMarker } = useDeleteMarker(
    currentMarkerInfo.markerId
  );

  const [isReview, setIsReview] = useState(false);

  const [isRoadView, setIsRoadView] = useState(false);

  const [isRoadViewError, setIsRoadViewError] = useState(false);
  const [roadViewErrorText, setRoadViewErrorText] = useState("");

  const [deleteLoading, setDeleteLoading] = useState(false);

  const [viewInput, setViewInput] = useState(false);

  useEffect(() => {
    toastState.close();
    toastState.setToastText("");
  }, []);

  useEffect(() => {
    if (!isRoadViewError) return;

    setRoadViewErrorText("거리뷰가 지원되지 않는 주소입니다");

    const errorTime = setTimeout(() => {
      setRoadViewErrorText("");
      setIsRoadViewError(false);
    }, 2000);

    return () => {
      clearTimeout(errorTime);
    };
  }, [isRoadViewError]);

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

  const handleLike = async () => {
    if (likePending || unLikePending) return;
    try {
      if (marker?.favorited) {
        await unLike();
      } else {
        await like();
      }
    } catch (error) {
      if (isAxiosError(error)) {
        if (error.response?.status === 401) {
          modalState.openLogin();
        }
      }
    }
  };

  const copyTextToClipboard = async () => {
    const url = `https://k-pullup.com/marker?d=${marker?.markerId}&la=${marker?.latitude}&lo=${marker?.longitude}`;
    try {
      await navigator.clipboard.writeText(url);
      toastState.setToastText("복사 완료");
      toastState.open();
    } catch (err) {
      alert("잠시 후 다시 시도해 주세요!");
    }
  };

  if (isLoading) return <MarkerInfoSkeleton />;
  if (isError) {
    if (isAxiosError(error)) {
      if (error.response?.status === 404) {
        console.log(error.response?.status);
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
            존재 하지 않는 장소 입니다...
          </div>
        );
      } else {
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
            잠시 후 다시 시도해 주세요...
          </div>
        );
      }
    }
  }

  return (
    <div>
      <Helmet>
        <title>대한민국 철봉 지도</title>
        <meta
          name="description"
          content={`${marker?.description} 위치: ${marker?.addr}`}
        />
        <meta
          property="og:image"
          content={marker?.photos ? marker.photos[0].photoUrl : noimg}
        />
      </Helmet>

      {isReview ? (
        <MarkerReview
          setIsReview={setIsReview}
          markerId={marker?.markerId as number}
        />
      ) : (
        <Styled.Container>
          <Tooltip title="닫기" arrow disableInteractive>
            <IconButton
              onClick={() => {
                setMarkerInfoModal(false);
                if (sharedMarker && sharedMarkerLat && sharedMarkerLng) {
                  navigate("/");
                }
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
            {isRoadView ? (
              <MarkerRoadView
                lat={marker?.latitude as number}
                lng={marker?.longitude as number}
                setIsRoadView={setIsRoadView}
                setIsRoadViewError={setIsRoadViewError}
              />
            ) : (
              <>
                <img
                  src={marker?.photos ? marker.photos[0].photoUrl : noimg}
                  alt="철봉 상세 이미지"
                />
                <Styled.description>
                  {marker?.isChulbong && (
                    <>
                      <Tooltip title="수정" arrow disableInteractive>
                        <IconButton
                          onClick={() => {
                            setViewInput(true);
                          }}
                          aria-label="edit"
                          sx={{
                            position: "absolute",
                            top: "0",
                            left: "0",
                            color: "#fff",
                          }}
                        >
                          <EditIcon
                            sx={{
                              width: "20px",
                              height: "20px",
                            }}
                          />
                        </IconButton>
                      </Tooltip>
                    </>
                  )}

                  {viewInput ? (
                    <Styled.InputWrap>
                      <Styled.DescInput
                        id="edit"
                        type="text"
                        onChange={descInput.onChange}
                      />
                      <Styled.ButtonWrap>
                        <ActionButton
                          bg="black"
                          onClick={() => {
                            updateDesc();
                            setViewInput(false);
                          }}
                        >
                          수정
                        </ActionButton>
                        <ActionButton
                          bg="gray"
                          onClick={() => {
                            setViewInput(false);
                          }}
                        >
                          취소
                        </ActionButton>
                      </Styled.ButtonWrap>
                    </Styled.InputWrap>
                  ) : (
                    <div>
                      {marker?.description || "작성된 설명이 없습니다."}
                    </div>
                  )}
                </Styled.description>
              </>
            )}
          </Styled.imageWrap>
          <Styled.AddressText>
            <div>
              {marker?.createdAt.toString().split("T")[0].replace(/-/g, ".")}{" "}
              등록
            </div>
            <div>
              {marker?.address ||
                "주소를 불러오는 중입니다. 잠시 후 다시 확인해 주세요!"}
            </div>
          </Styled.AddressText>

          <Styled.ErrorBox>{roadViewErrorText}</Styled.ErrorBox>
          <Styled.BottomButtons>
            <Tooltip title="리뷰 보기" arrow disableInteractive>
              <IconButton onClick={handleViewReview} aria-label="review">
                <RateReviewIcon />
              </IconButton>
            </Tooltip>
            <Tooltip title="거리뷰" arrow disableInteractive>
              <IconButton
                onClick={() => {
                  setIsRoadView(true);
                }}
                aria-label="review"
              >
                <PlaceIcon />
              </IconButton>
            </Tooltip>
            <Tooltip title="공유 링크 복사" arrow disableInteractive>
              <IconButton onClick={copyTextToClipboard} aria-label="dislike">
                <div
                  style={{
                    width: "24px",
                    height: "24px",
                    position: "relative",
                  }}
                >
                  <ShareIcon />
                </div>
              </IconButton>
            </Tooltip>
            {marker?.favorited ? (
              <Tooltip title="좋아요 취소" arrow disableInteractive>
                <IconButton onClick={handleLike} aria-label="dislike">
                  <div
                    style={{
                      width: "24px",
                      height: "24px",
                      position: "relative",
                    }}
                  >
                    <FavoriteIcon />
                  </div>
                </IconButton>
              </Tooltip>
            ) : (
              <Tooltip title="좋아요" arrow disableInteractive>
                <IconButton onClick={handleLike} aria-label="dislike">
                  <div
                    style={{
                      width: "24px",
                      height: "24px",
                      position: "relative",
                    }}
                  >
                    <FavoriteBorderIcon />
                  </div>
                </IconButton>
              </Tooltip>
            )}
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

            {marker?.isChulbong && (
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
        </Styled.Container>
      )}
    </div>
  );
};

export default MarkerInfoModal;
