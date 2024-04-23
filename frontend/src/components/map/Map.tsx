"use client";

import useAllMarkerData from "@/hooks/query/useAllMarkerData";
import useBodyToggleStore from "@/store/useBodyToggleStore";
import useMapStatusStore from "@/store/useMapStatusStore";
import useMapStore from "@/store/useMapStore";
import { useEffect, useRef, useState } from "react";
import MapLoading from "./MapLoading";
import getWeather from "@/api/markers/getWeather";
import getMarker from "@/api/markers/getMarker";
import useMobileMapOpenStore from "@/store/useMobileMapOpenStore";

const Map = () => {
  const { isOpen: isMobileMapOpen } = useMobileMapOpenStore();
  const { lat, lng, level, setLevel, setPosition } = useMapStatusStore();

  const { map, setMap, setClusterer } = useMapStore();
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

      window.kakao.maps.event.addListener(newMarker, "click", async () => {
        content.innerHTML = "";
        const infoBox = `
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
              <a>상세보기</a>
              <a>정보 수정 제안</a>
            </div>
            <div class="empty-grow"></div>
            <div id="overlay-action">
              <button>
                <div>
                  <img src="/bookmark-02.svg" alt="bookmark"/>
                </div>
                <div>
                  저장
                </div>
              </button>
              <button>
                <div>
                  <img src="/roadview.svg" alt="roadview"/>
                </div>
                <div>
                  거리뷰
                </div>
              </button>
              <button>
                <div>
                  <img src="/share-08.svg" alt="share"/>
                </div>
                <div>
                  거리뷰
                </div>
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

        const { iconImage, temperature, desc } = await getWeather(
          marker.latitude,
          marker.longitude
        );
        const { description, address, favorited, photos } = await getMarker(
          marker.markerId
        );

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
        imageBox.src = photos ? photos[0].photoUrl : "/metaimg.webp";
        imageBox.onload = () => {
          imageBox.style.display = "block";
          imageContainer.classList.remove("on-loading");
        };

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
