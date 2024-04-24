"use client";

import useAllMarkerData from "@/hooks/query/marker/useAllMarkerData";
import useBodyToggleStore from "@/store/useBodyToggleStore";
import useMapStatusStore from "@/store/useMapStatusStore";
import useMapStore from "@/store/useMapStore";
import { useEffect, useRef, useState } from "react";
import MapLoading from "./MapLoading";
import getWeather from "@/api/markers/getWeather";
import getMarker from "@/api/markers/getMarker";
import useMobileMapOpenStore from "@/store/useMobileMapOpenStore";
import { useRouter } from "next/navigation";
import deleteFavorite from "@/api/favorite/deleteFavorite";
import setFavorite from "@/api/favorite/setFavorite";
import { type Photo } from "@/types/Marker.types";

const Map = () => {
  const router = useRouter();

  const { isOpen: isMobileMapOpen } = useMobileMapOpenStore();
  const { lat, lng, level, setLevel, setPosition } = useMapStatusStore();

  const { map, setMap, setClusterer, setMarkers } = useMapStore();
  const { isOpen } = useBodyToggleStore();

  const { data: markers } = useAllMarkerData();

  const [mapLoading, setMapLoading] = useState(true);

  const mapRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (!markers) return;

    const options = {
      center: new window.kakao.maps.LatLng(lat, lng),
      level: level,
      maxLevel: 12,
    };

    const newMap = new window.kakao.maps.Map(mapRef.current, options);

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

    const imageSize = new window.kakao.maps.Size(39, 39);
    const imageOption = { offset: new window.kakao.maps.Point(27, 45) };

    const activeMarkerImg = new window.kakao.maps.MarkerImage(
      "/activeMarker.svg",
      imageSize,
      imageOption
    );

    const clusterer = new window.kakao.maps.MarkerClusterer({
      map: newMap,
      averageCenter: true,
      minLevel: 6,
    });

    const skeletoncontent = document.createElement("div");
    skeletoncontent.className = "skeleton-overlay";

    const content = document.createElement("div");

    const skeletonOverlay = new window.kakao.maps.CustomOverlay({
      content: skeletoncontent,
      zIndex: 5,
    });

    const newMarkers = markers?.map((marker) => {
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
        if (weatherLoading || markerLoading) return;
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
                <button>정보 수정 제안</button>
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
                <button>
                  <div>
                    <img src="/roadview.svg" alt="roadview" />
                  </div>
                  <div>거리뷰</div>
                </button>
                <button>
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

        const latlng = new window.kakao.maps.LatLng(
          marker.latitude,
          marker.longitude
        );

        skeletonOverlay.setMap(newMap);
        skeletonOverlay.setPosition(latlng);

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
            addBookmarkError = true;
          } finally {
            addBookmarkLoading = false;
          }
        };

        const deleteBookmark = async () => {
          deleteBookmarkLoading = true;
          try {
            const res = await deleteFavorite(marker.markerId);
            return res;
          } catch (error) {
            deleteBookmarkError = true;
          } finally {
            deleteBookmarkLoading = false;
          }
        };

        await fetchMarker();
        await fetchWeather();

        skeletonOverlay.setMap(null);

        overlay.setMap(newMap);
        overlay.setPosition(latlng);

        // 오버레이 날씨 정보
        const weatherIconBox = document.getElementById(
          "overlay-weather-icon"
        ) as HTMLImageElement;
        weatherIconBox.src = `${iconImage}` || "";
        weatherIconBox.alt = `${desc} || ""`;
        const weatherTempBox = document.getElementById(
          "overlay-weather-temp"
        ) as HTMLDivElement;
        weatherTempBox.innerHTML = `${temperature}℃`;

        // 오버레이 주소 정보
        const addressBox = document.getElementById(
          "overlay-title"
        ) as HTMLDivElement;
        addressBox.innerHTML = description || "작성된 설명이 없습니다.";

        // 오버레이 이미지 정보
        const imageContainer = document.getElementById(
          "overlay-image-container"
        ) as HTMLDivElement;
        imageContainer.classList.add("on-loading");
        const imageBox = document.getElementById(
          "overlay-image"
        ) as HTMLImageElement;
        imageBox.src = photos ? photos[0]?.photoUrl : "/metaimg.webp";
        imageBox.onload = () => {
          imageBox.style.display = "block";
          imageContainer.classList.remove("on-loading");
        };

        // 오버레이 상세보기 링크
        const detailLink = document.getElementById(
          "item-detail-link"
        ) as HTMLAnchorElement;
        detailLink.style.cursor = "pointer";
        detailLink.addEventListener("click", () => {
          router.push(`/pullup/${marker.markerId}`);
        });

        // 오버레이 북마크 버튼 이미지
        const bookmarkBtnImg = document.getElementById(
          "bookmark-button-img"
        ) as HTMLImageElement;
        bookmarkBtnImg.src = favorited
          ? "/bookmark-03.svg"
          : "/bookmark-02.svg";

        // 오버레이 북마크 버튼 액션
        const bookmarkBtn = document.getElementById(
          "bookmark-button-img"
        ) as HTMLButtonElement;
        const bookmarkText = document.getElementById(
          "bookmark-text"
        ) as HTMLDivElement;
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

        // 오버레이 닫기 이벤트 등록
        const closeBtnBox = document.getElementById(
          "overlay-close"
        ) as HTMLButtonElement;
        closeBtnBox.onclick = () => {
          overlay.setMap(null);
        };
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
    if (!map) return;
    const moveLatLon = new window.kakao.maps.LatLng(lat, lng);

    map.relayout();

    const resizeTime = setTimeout(() => {
      map.setCenter(moveLatLon);
      map.relayout();
    }, 200);

    return () => clearTimeout(resizeTime);
  }, [isOpen, mapLoading, map, isMobileMapOpen]);

  return (
    <div className="relative w-full h-screen">
      {mapLoading && <MapLoading />}
      <div
        id="map"
        ref={mapRef}
        className={`absolute top-0 left-0 w-full h-full ${
          mapLoading ? "hidden" : "block"
        }`}
      />
    </div>
  );
};

export default Map;
