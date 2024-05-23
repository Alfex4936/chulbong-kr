"use client";

import deleteFavorite from "@/api/favorite/deleteFavorite";
import setFavorite from "@/api/favorite/setFavorite";
import getMarker from "@/api/markers/getMarker";
import getWeather from "@/api/markers/getWeather";
import { MOBILE_WIDTH } from "@/constants";
import useAllMarkerData from "@/hooks/query/marker/useAllMarkerData";
import useBodyToggleStore from "@/store/useBodyToggleStore";
import useLoginModalStateStore from "@/store/useLoginModalStateStore";
import useMapStatusStore from "@/store/useMapStatusStore";
import useMapStore from "@/store/useMapStore";
import useMobileMapOpenStore from "@/store/useMobileMapOpenStore";
import usePageLoadingStore from "@/store/usePageLoadingStore";
import useRoadviewStatusStore from "@/store/useRoadviewStatusStore";
import useSelectedMarkerStore from "@/store/useSelectedMarkerStore";
import { type Photo } from "@/types/Marker.types";
import { useQueryClient } from "@tanstack/react-query";
import { isAxiosError } from "axios";
import { usePathname, useRouter } from "next/navigation";
import { useEffect, useRef, useState } from "react";
import { useToast } from "../ui/use-toast";
import MapLoading from "./MapLoading";
import MapSearch from "./MapSearch";
import { MdOutlineGpsFixed } from "react-icons/md";
import MapButtons from "./MapButtons";
import { type CustomOverlay } from "@/types/CustomOverlay.types";
import MyLocateOverlay from "./MyLocateOverlay";
import { createRoot } from "react-dom/client";
import PlusIcon from "../icons/PlusIcon";
import MinusIcon from "../icons/MinusIcon";

const Map = () => {
  const pathname = usePathname();
  const router = useRouter();

  const queryClient = useQueryClient();

  const path1 = pathname.split("/")[1];

  const { open } = useLoginModalStateStore();

  const { marker: selectedMarker } = useSelectedMarkerStore();

  const { close: mobileMapClose } = useMobileMapOpenStore();

  const { isOpen: isMobileMapOpen } = useMobileMapOpenStore();
  const { lat, lng, level, setLevel, setPosition } = useMapStatusStore();
  const { open: openRoadview, setPosition: setRoadview } =
    useRoadviewStatusStore();

  const { toast } = useToast();

  const {
    map,
    setMap,
    clusterer: sClusterer,
    setClusterer,
    setMarkers,
    setOverlay,
    markers: mapMarkers,
    overlay: overlayState,
  } = useMapStore();
  const { isOpen } = useBodyToggleStore();

  const { setLoading } = usePageLoadingStore();

  const { data: markers } = useAllMarkerData();

  const [mapLoading, setMapLoading] = useState(true);

  const [bookmarkError, setBookmarkError] = useState(false);

  const [currentOverlay, setCurrentOverlay] = useState<CustomOverlay | null>(
    null
  ); // GPS 현재 위치

  const mapRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (!markers) return;

    const imageSize = new window.kakao.maps.Size(39, 39);
    const imageOption = { offset: new window.kakao.maps.Point(27, 45) };

    const activeMarkerImg = new window.kakao.maps.MarkerImage(
      "/activeMarker.svg",
      imageSize,
      imageOption
    );

    const options = {
      center: new window.kakao.maps.LatLng(lat, lng),
      level: level,
      maxLevel: 12,
    };

    const newMap = map || new window.kakao.maps.Map(mapRef.current, options);

    const handleDrag = () => {
      const latlng = newMap.getCenter();
      setPosition(latlng.getLat(), latlng.getLng());
    };

    const handleZoom = () => {
      const level = newMap.getLevel();
      setLevel(level);
    };

    window.kakao.maps.event.addListener(newMap, "dragend", handleDrag);
    window.kakao.maps.event.addListener(newMap, "zoom_changed", handleZoom);

    const clusterer =
      sClusterer ||
      new window.kakao.maps.MarkerClusterer({
        map: newMap,
        averageCenter: true,
        minLevel: 6,
      });

    if (sClusterer) sClusterer.clear();

    const skeletoncontent = document.createElement("div");
    skeletoncontent.className = "skeleton-overlay";

    const content = document.createElement("div");

    const skeletonOverlay = new window.kakao.maps.CustomOverlay({
      content: skeletoncontent,
      zIndex: 5,
    });

    const newMarkers = markers?.map((marker) => {
      const changeRoadviewlocation = async () => {
        setRoadview(marker.latitude, marker.longitude);
      };

      const copyTextToClipboard = async () => {
        const url = `${process.env.NEXT_PUBLIC_URL}/pullup/${marker.markerId}`;
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
        position: new window.kakao.maps.LatLng(
          marker.latitude,
          marker.longitude
        ),
        image: activeMarkerImg,
        title: marker.markerId,
        zIndex: 4,
      });

      let markerLoading = false;
      let weatherLoading = false;

      window.kakao.maps.event.addListener(newMarker, "click", async () => {
        if (document.getElementsByClassName("overlay")[0]) {
          document.getElementsByClassName("overlay")[0].remove();
        }

        if (weatherLoading || markerLoading) return;
        const latlng = new window.kakao.maps.LatLng(
          marker.latitude,
          marker.longitude
        );

        skeletonOverlay.setMap(newMap);
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
            const res = await getMarker(marker.markerId);
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
            const res = await getWeather(marker.latitude, marker.longitude);
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
            const res = await setFavorite(marker.markerId);
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
              queryKey: ["marker", marker.markerId],
            });
          }
        };

        const deleteBookmark = async () => {
          deleteBookmarkLoading = true;
          try {
            const res = await deleteFavorite(marker.markerId);
            return res;
          } catch (error) {
            deleteBookmarkError = true;
            setBookmarkError(true);
          } finally {
            deleteBookmarkLoading = false;
            queryClient.invalidateQueries({
              queryKey: ["marker", marker.markerId],
            });
          }
        };

        await fetchMarker();
        await fetchWeather();

        skeletonOverlay.setMap(null);

        overlay.setMap(newMap);
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
            setLoading(true);
            if (window.innerWidth <= MOBILE_WIDTH) {
              mobileMapClose();
            }
            router.push(`/pullup/${marker.markerId}`);
          });
        }

        // 오버레이 정보 수정 제안 요청
        const reportLink = document.getElementById(
          "item-report"
        ) as HTMLAnchorElement;
        if (reportLink) {
          reportLink.style.cursor = "pointer";
          reportLink.addEventListener("click", () => {
            setLoading(true);
            if (window.innerWidth <= MOBILE_WIDTH) {
              mobileMapClose();
            }
            router.push(`/pullup/${marker.markerId}/reportlist`);
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

      return newMarker;
    });

    clusterer.addMarkers(newMarkers);

    setMapLoading(false);
    setMap(newMap);
    setClusterer(clusterer);
    setMarkers(newMarkers);
    return () => {
      window.kakao.maps.event.removeListener(newMap, "dragend", handleDrag);
      window.kakao.maps.event.removeListener(
        newMap,
        "zoom_changed",
        handleZoom
      );
    };
  }, [markers]);

  useEffect(() => {
    if (!map || !mapMarkers) return;

    let resizeTime: NodeJS.Timeout;

    if (path1 === "pullup") {
      if (!selectedMarker) return;

      const filterClickMarker = () => {
        const imageSize = new window.kakao.maps.Size(39, 39);
        const imageOption = { offset: new window.kakao.maps.Point(27, 45) };
        const selectedMarkerImg = new window.kakao.maps.MarkerImage(
          "/selectedMarker.svg",
          imageSize,
          imageOption
        );
        const activeMarkerImg = new window.kakao.maps.MarkerImage(
          "/activeMarker.svg",
          imageSize,
          imageOption
        );
        mapMarkers.forEach((marker) => {
          if (Number(marker.getTitle()) === selectedMarker.markerId) {
            marker.setImage(selectedMarkerImg);
          } else {
            marker.setImage(activeMarkerImg);
          }
        });
      };

      const moveLatLon = new window.kakao.maps.LatLng(
        selectedMarker.lat,
        selectedMarker.lng
      );

      map.relayout();

      resizeTime = setTimeout(() => {
        map.setCenter(moveLatLon);
        map.relayout();
        filterClickMarker();
      }, 200);
    } else {
      const moveLatLon = new window.kakao.maps.LatLng(lat, lng);

      map.relayout();

      resizeTime = setTimeout(() => {
        map.setCenter(moveLatLon);
        map.relayout();
      }, 200);
    }
    return () => clearTimeout(resizeTime);
  }, [mapLoading, map, isMobileMapOpen, path1, selectedMarker]);

  useEffect(() => {
    if (!map) return;

    const resizeTime = setTimeout(() => {
      map.relayout();
    }, 200);

    return () => clearTimeout(resizeTime);
  }, [isOpen]);

  useEffect(() => {
    if (!map) return;
    const moveLatLon = new window.kakao.maps.LatLng(lat, lng);

    map.relayout();
    map.setCenter(moveLatLon);
  }, [map]);

  const centerMapOnCurrentPosition = () => {
    if (map && navigator.geolocation) {
      navigator.geolocation.getCurrentPosition(
        (position) => {
          const moveLatLon = new window.kakao.maps.LatLng(
            position.coords.latitude,
            position.coords.longitude
          );

          if (currentOverlay) {
            currentOverlay.setMap(null);
          }

          const overlayDiv = document.createElement("div");
          const root = createRoot(overlayDiv);
          root.render(<MyLocateOverlay />);

          const customOverlay = new window.kakao.maps.CustomOverlay({
            position: moveLatLon,
            content: overlayDiv,
            zIndex: 3,
          });

          customOverlay.setMap(map);
          setCurrentOverlay(customOverlay);

          setPosition(position.coords.latitude, position.coords.longitude);
          map.setCenter(moveLatLon);
        },
        () => {
          alert("잠시 후 다시 시도해주세요.");
        }
      );
    } else {
      alert("잠시 후 다시 시도해주세요.");
    }
  };

  const zoomIn = () => {
    const level = map?.getLevel();

    map?.setLevel((level as number) - 1);
  };

  const zoomOut = () => {
    const level = map?.getLevel();

    map?.setLevel((level as number) + 1);
  };

  return (
    <div className="relative w-full h-dvh">
      {mapLoading && <MapLoading />}
      <div
        id="map"
        ref={mapRef}
        className={`absolute top-0 left-0 w-full h-full ${
          mapLoading ? "hidden" : "block"
        }`}
      >
        <MapSearch />
        <MapButtons
          icon={<MdOutlineGpsFixed />}
          className="top-16 right-2"
          onClick={centerMapOnCurrentPosition}
          tooltipText="내 위치"
        />
        <MapButtons
          icon={<PlusIcon size={18} />}
          className="top-28 right-2"
          onClick={zoomIn}
          tooltipText="확대"
        />
        <MapButtons
          icon={<MinusIcon size={18} />}
          className="top-[143px] right-2"
          onClick={zoomOut}
          tooltipText="축소"
        />
      </div>
    </div>
  );
};

export default Map;
