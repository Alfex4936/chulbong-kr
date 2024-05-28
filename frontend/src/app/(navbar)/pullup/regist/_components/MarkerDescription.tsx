"use client";

import deleteFavorite from "@/api/favorite/deleteFavorite";
import setFavorite from "@/api/favorite/setFavorite";
import getMarker from "@/api/markers/getMarker";
import getWeather from "@/api/markers/getWeather";
import ErrorMessage from "@/components/atom/ErrorMessage";
import LoadingSpinner from "@/components/atom/LoadingSpinner";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { useToast } from "@/components/ui/use-toast";
import { MOBILE_WIDTH } from "@/constants";
import useInput from "@/hooks/common/useInput";
import useMapControl from "@/hooks/common/useMapControl";
import useSetFacilities from "@/hooks/mutation/marker/useSetFacilities";
import useUploadMarker from "@/hooks/mutation/marker/useUploadMarker";
import useReportMarker from "@/hooks/mutation/report/useReportMarker";
import useMapStore from "@/store/useMapStore";
import useMobileMapOpenStore from "@/store/useMobileMapOpenStore";
import useRoadviewStatusStore from "@/store/useRoadviewStatusStore";
import useUploadFormDataStore from "@/store/useUploadFormDataStore";
import { Photo } from "@/types/Marker.types";
import { useQueryClient } from "@tanstack/react-query";
import { isAxiosError } from "axios";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";

interface Props {
  desc?: string;
  markerId?: number;
  isReport?: boolean;
}

const MarkerDescription = ({ desc, markerId, isReport = false }: Props) => {
  const descriptionValue = useInput(desc ? desc : "");

  const queryClient = useQueryClient();

  const router = useRouter();

  const { close: mobileMapClose } = useMobileMapOpenStore();
  const { toast } = useToast();
  const { open: openRoadview, setPosition: setRoadview } =
    useRoadviewStatusStore();

  const { mutateAsync: setFacilities } = useSetFacilities();
  const { mutateAsync: uploadMarker } = useUploadMarker();
  const {
    clusterer,
    map,
    setMarkers,
    markers,
    setOverlay,
    overlay: overlayState,
  } = useMapStore();
  const { filterMarker, moveLocation } = useMapControl();

  const [loading, setLoading] = useState(false);
  const [bookmarkError, setBookmarkError] = useState(false);
  const [error, setError] = useState("");

  const [errorMessage, setErrorMessage] = useState("");

  const { imageForm, latitude, longitude, resetData, facilities } =
    useUploadFormDataStore();

  const {
    mutate: report,
    error: reportError,
    isError,
  } = useReportMarker(markerId as number);

  useEffect(() => {
    resetData();
  }, []);

  const handleSubmit = async () => {
    if (isReport && markerId) {
      setLoading(true);
      try {
        const marker = await getMarker(markerId);
        let data;
        if (latitude !== 0 && longitude !== 0) {
          data = {
            markerId: markerId,
            description: descriptionValue.value,
            photos: imageForm.map((image) => image.file) as File[],
            latitude: marker.latitude,
            longitude: marker.longitude,
            newLatitude: latitude,
            newLongitude: longitude,
          };
        } else {
          data = {
            markerId: markerId,
            description: descriptionValue.value,
            photos: imageForm.map((image) => image.file) as File[],
            latitude: marker.latitude,
            longitude: marker.longitude,
          };
        }

        if (imageForm.length <= 0) {
          setErrorMessage("이미지를 1개 이상 등록해 주세요");
          setLoading(false);
          return;
        }

        report(data);
      } catch (error) {
        setErrorMessage("잠시 후 다시 시도해 주세요.");
        setLoading(false);
      }

      return;
    }

    const data = {
      description: descriptionValue.value,
      photos: imageForm.map((image) => image.file) as File[],
      latitude: latitude,
      longitude: longitude,
    };

    setLoading(true);

    try {
      const result = await uploadMarker(data);
      await setFacilities({
        markerId: result.markerId,
        facilities: [
          {
            facilityId: 1,
            quantity: facilities.철봉,
          },
          {
            facilityId: 2,
            quantity: facilities.평행봉,
          },
        ],
      });

      const imageSize = new window.kakao.maps.Size(39, 39);
      const imageOption = { offset: new window.kakao.maps.Point(27, 45) };

      const activeMarkerImg = new window.kakao.maps.MarkerImage(
        "/selectedMarker.svg",
        imageSize,
        imageOption
      );

      const skeletoncontent = document.createElement("div");
      skeletoncontent.className = "skeleton-overlay";

      const content = document.createElement("div");

      const skeletonOverlay = new window.kakao.maps.CustomOverlay({
        content: skeletoncontent,
        zIndex: 5,
      });

      const changeRoadviewlocation = async () => {
        setRoadview(latitude, longitude);
      };

      const copyTextToClipboard = async () => {
        const url = `${process.env.NEXT_PUBLIC_URL}/pullup/${result.markerId}`;
        try {
          await navigator.clipboard.writeText(url);
          toast({
            description: "링크 복사 완료",
          });
        } catch (err) {
          alert("잠시 후 다시 시도해 주세요!");
        }
      };

      const newMarker = new window.kakao.maps.Marker({
        position: new window.kakao.maps.LatLng(latitude, longitude),
        image: activeMarkerImg,
        title: result.markerId,
        zIndex: 4,
      });

      let markerLoading = false;
      let weatherLoading = false;

      window.kakao.maps.event.addListener(newMarker, "click", async () => {
        if (document.getElementsByClassName("overlay")[0]) {
          document.getElementsByClassName("overlay")[0].remove();
        }

        if (weatherLoading || markerLoading) return;

        const latlng = new window.kakao.maps.LatLng(latitude, longitude);

        skeletonOverlay.setMap(map);
        skeletonOverlay.setPosition(latlng);

        content.innerHTML = "";
        const infoBox = /* HTML */ `
          <div id="overlay-top">
            <div id="overlay-weather">
              <div>
                <img id="overlay-weather-icon" />
              </div>
              <div id="overlay-weather-temp"></div>
            </div>
            <button id="overlay-close">닫기</button>
          </div>
          <div id="overlay-mid">
            <div id="overlay-info">
              <div id="overlay-title"></div>
              <div id="overlay-link">
                <button id="item-detail-link">상세보기</button>
                <button id="item-report">정보 수정 제안</button>
              </div>
              <div class="empty-grow"></div>
              <div id="overlay-action">
                <button id="bookmark-button">
                  <div>
                    <img
                      id="bookmark-button-img"
                      src="/bookmark-02.svg"
                      alt="bookmark"
                    />
                  </div>
                  <div id="bookmark-text">북마크</div>
                </button>
                <button id="roadview-button">
                  <div>
                    <img src="/roadview.svg" alt="roadview" />
                  </div>
                  <div>거리뷰</div>
                </button>
                <button id="share-button">
                  <div>
                    <img src="/share-08.svg" alt="share" />
                  </div>
                  <div>공유</div>
                </button>
              </div>
            </div>
            <div id="overlay-image-container">
              <img id="overlay-image" />
            </div>
          </div>
        `;

        content.className = "overlay";
        content.innerHTML = infoBox;

        const overlay = new window.kakao.maps.CustomOverlay({
          content: content,
          zIndex: 5,
        });

        setOverlay(overlay);

        // 마커 정보
        let description: string = "";
        let address: string = "";
        let favorited: boolean = false;
        let photos: Photo[] = [];
        let markerError = false;
        // 날씨 정보
        let iconImage: string = "";
        let temperature: string = "";
        let desc: string = "";
        let weatherError = false;
        // 북마크 정보
        let addBookmarkLoading = false;
        let addBookmarkError = false;
        let deleteBookmarkLoading = false;
        let deleteBookmarkError = false;

        const fetchMarker = async () => {
          markerLoading = true;
          try {
            const res = await getMarker(result.markerId);
            description = res.description;
            address = res.address as string;
            favorited = res.favorited as boolean;
            photos = res.photos as Photo[];
          } catch (error) {
            markerError = true;
            content.innerHTML = /* HTML */ `
              <div class="error-box">
                <span>잘못된 위치입니다. 잠시 후 다시 시도해 주세요.</span>
                <span><button id="error-close">닫기</button></span>
              </div>
            `;
            const errorCloseBtn = document.getElementById("error-close");
            errorCloseBtn?.addEventListener("click", () => {
              overlay.setMap(null);
            });
          } finally {
            markerLoading = false;
          }
        };

        const fetchWeather = async () => {
          weatherLoading = true;
          try {
            const res = await getWeather(latitude, longitude);
            iconImage = res.iconImage;
            temperature = res.temperature;
            desc = res.desc;
          } catch (error) {
            weatherError = true;
          } finally {
            weatherLoading = false;
          }
        };

        const addBookmark = async () => {
          addBookmarkLoading = true;
          try {
            const res = await setFavorite(result.markerId);
            return res;
          } catch (error) {
            if (isAxiosError(error)) {
              if (error.response?.status === 401) open();
            } else {
              toast({
                description: "잠시 후 다시 시도해 주세요",
              });
            }
            addBookmarkError = true;
            setBookmarkError(true);
          } finally {
            addBookmarkLoading = false;
            queryClient.invalidateQueries({
              queryKey: ["marker", result.markerId],
            });
          }
        };

        const deleteBookmark = async () => {
          deleteBookmarkLoading = true;
          try {
            const res = await deleteFavorite(result.markerId);
            return res;
          } catch (error) {
            deleteBookmarkError = true;
            setBookmarkError(true);
          } finally {
            deleteBookmarkLoading = false;
            queryClient.invalidateQueries({
              queryKey: ["marker", result.markerId],
            });
          }
        };

        await fetchMarker();
        await fetchWeather();

        skeletonOverlay.setMap(null);

        overlay.setMap(map);
        overlay.setPosition(latlng);

        setOverlay(overlay);

        // 오버레이 날씨 정보
        const weatherIconBox = document.getElementById(
          "overlay-weather-icon"
        ) as HTMLImageElement;
        if (weatherIconBox) {
          weatherIconBox.src = `${iconImage}` || "";
          weatherIconBox.alt = `${desc} || ""`;
        }

        const weatherTempBox = document.getElementById(
          "overlay-weather-temp"
        ) as HTMLDivElement;
        if (weatherTempBox) {
          weatherTempBox.innerHTML = `${temperature}℃`;
        }

        // 오버레이 주소 정보
        const addressBox = document.getElementById(
          "overlay-title"
        ) as HTMLDivElement;
        if (addressBox) {
          addressBox.innerHTML = description || "작성된 설명이 없습니다.";
        }

        // 오버레이 이미지 정보
        const imageContainer = document.getElementById(
          "overlay-image-container"
        ) as HTMLDivElement;
        if (imageContainer) {
          imageContainer.classList.add("on-loading");
        }
        const imageBox = document.getElementById(
          "overlay-image"
        ) as HTMLImageElement;
        if (imageBox) {
          imageBox.src = photos ? photos[0]?.photoUrl : "/metaimg.webp";
          imageBox.onload = () => {
            imageBox.style.display = "block";
            imageContainer.classList.remove("on-loading");
          };
        }

        // 오버레이 상세보기 링크
        const detailLink = document.getElementById(
          "item-detail-link"
        ) as HTMLAnchorElement;
        if (detailLink) {
          detailLink.style.cursor = "pointer";
          detailLink.addEventListener("click", () => {
            if (window.innerWidth <= MOBILE_WIDTH) {
              mobileMapClose();
            }
            router.push(`/pullup/${result.markerId}`);
          });
        }

        // 오버레이 정보 수정 제안 요청
        const reportLink = document.getElementById(
          "item-report"
        ) as HTMLAnchorElement;
        if (reportLink) {
          reportLink.style.cursor = "pointer";
          reportLink.addEventListener("click", () => {
            if (window.innerWidth <= MOBILE_WIDTH) {
              mobileMapClose();
            }
            router.push(`/pullup/${result.markerId}/reportlist`);
          });
        }

        // 오버레이 북마크 버튼 이미지
        const bookmarkBtnImg = document.getElementById(
          "bookmark-button-img"
        ) as HTMLImageElement;
        if (bookmarkBtnImg) {
          bookmarkBtnImg.src = favorited
            ? "/bookmark-03.svg"
            : "/bookmark-02.svg";
        }

        // 오버레이 북마크 버튼 액션
        const bookmarkBtn = document.getElementById(
          "bookmark-button-img"
        ) as HTMLButtonElement;
        const bookmarkText = document.getElementById(
          "bookmark-text"
        ) as HTMLDivElement;
        if (bookmarkBtn && bookmarkText) {
          bookmarkBtn.addEventListener("click", async () => {
            if (addBookmarkLoading || deleteBookmarkLoading) return;
            bookmarkBtn.disabled = true;
            if (favorited) {
              bookmarkText.innerHTML = "취소중..";
              await deleteBookmark();
            } else if (!favorited) {
              bookmarkText.innerHTML = "저장중..";
              await addBookmark();
            }
            await fetchMarker();

            bookmarkText.innerHTML = "북마크";
            bookmarkBtnImg.src = favorited
              ? "/bookmark-03.svg"
              : "/bookmark-02.svg";

            bookmarkBtn.disabled = false;
          });
        }

        // 오보레이 로드뷰 버튼
        const roadviewButton = document.getElementById(
          "roadview-button"
        ) as HTMLButtonElement;
        if (roadviewButton) {
          roadviewButton.addEventListener("click", async () => {
            await changeRoadviewlocation();
            openRoadview();
          });
        }

        // 오버레이 공유 버튼
        const shareButton = document.getElementById(
          "share-button"
        ) as HTMLButtonElement;
        if (shareButton) {
          shareButton.addEventListener("click", copyTextToClipboard);
        }

        // 오버레이 닫기 이벤트 등록
        const closeBtnBox = document.getElementById(
          "overlay-close"
        ) as HTMLButtonElement;
        if (closeBtnBox) {
          closeBtnBox.onclick = () => {
            overlay.setMap(null);
          };
        }

        // 에러 오버레이 닫기
        const errorCloseBtn = document.getElementById("error-close");
        if (errorCloseBtn) {
          errorCloseBtn.onclick = () => {
            overlay.setMap(null);
          };
        }
      });

      if (!markers) return;
      const newMarkers = [...markers, newMarker];
      clusterer?.addMarker(newMarker);

      setMarkers(newMarkers);

      await filterMarker(result.markerId);
      await moveLocation(latitude, longitude);

      router.push(`/pullup/${result.markerId}`);
    } catch (error) {
      if (isAxiosError(error)) {
        if (error.response?.status === 401) {
        } else if (error.response?.status === 409) {
          setError("주변에 이미 철봉이 있습니다!");
        } else if (error.response?.status === 403) {
          setError("대한민국에서만 등록 가능합니다!");
        } else if (error.response?.status === 400) {
          setError("입력을 확인해 주세요!");
        } else {
          setError("잠시 후 다시 시도해 주세요!");
        }
      }
      setLoading(false);
    }
  };

  useEffect(() => {
    if (isAxiosError(reportError)) {
      if (reportError.response?.status === 406) {
        setErrorMessage("기존 위치와 너무 많은 거리가 차이납니다.");
      }
    }
  }, [isError]);

  return (
    <div>
      <div className="flex flex-col mb-5">
        <Input
          className="text-base"
          type="text"
          id="description"
          placeholder={desc ? "" : "설명 입력"}
          maxLength={70}
          value={descriptionValue.value}
          onChange={(e) => {
            if (descriptionValue.value.length >= 70) {
              setError("70자 이내로 작성해 주세요!");
            } else {
              setError("");
            }
            descriptionValue.handleChange(e);
          }}
        />
        <div
          data-testid="file-error"
          className="mt-1 text-center text-sm text-red"
        >
          {error}
        </div>
      </div>

      <div>
        <Button
          className="border-grey border bg-transparent dark:text-grey hover:bg-white-tp-light hover:border-transparent"
          onClick={handleSubmit}
          disabled={loading}
        >
          {loading ? (
            <LoadingSpinner size="xs" />
          ) : markerId && isReport ? (
            "제안 요청"
          ) : (
            "등록하기"
          )}
        </Button>
        {loading && imageForm.length > 0 ? (
          <div
            data-testid="file-error"
            className="mt-1 text-center text-sm text-red"
          >
            이미지를 등록하는 중입니다. 잠시만 기다려 주세요!
          </div>
        ) : null}
      </div>
      <ErrorMessage text={errorMessage} />
    </div>
  );
};

export default MarkerDescription;
