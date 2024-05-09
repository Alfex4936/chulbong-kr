"use client";

import useMapStatusStore from "@/store/useMapStatusStore";
import useMiniMapStatusStore from "@/store/useMiniMapStatusStore";
import useUploadFormDataStore from "@/store/useUploadFormDataStore";
import type { KaKaoMapMouseEvent } from "@/types/KakaoMap.types";
import getAddress, { type AddressInfo } from "@/utils/getAddress";
import { useCallback, useEffect, useRef, useState } from "react";
// 모바일 스크롤 멈춤 적용 안됨

interface Props {
  isMarker?: boolean;
  latitude?: number;
  longitude?: number;
}

const MiniMap = ({ isMarker = false, latitude, longitude }: Props) => {
  const { lat, lng, level } = useMapStatusStore();
  const { loaded, setMap, map } = useMiniMapStatusStore();
  const { setPosition: setFormPosition } = useUploadFormDataStore();

  const [mapOver, setMapOver] = useState(false);
  const [isClickMap, setIsClickMap] = useState(false);

  const [address, setAddress] = useState("");

  const mapRef = useRef<HTMLDivElement>(null);

  const setDraggable = useCallback(
    (draggable: boolean) => {
      // 마우스 드래그로 지도 이동 가능여부를 설정합니다
      map?.setDraggable(draggable);
    },
    [map]
  );

  useEffect(() => {
    if (!loaded) return;

    const options = {
      center: new window.kakao.maps.LatLng(lat, lng),
      level: level,
      maxLevel: 12,
    };

    const newMap = new window.kakao.maps.Map(mapRef.current, options);

    newMap.setDraggable(false);

    setMap(newMap);

    const imageSize = new window.kakao.maps.Size(39, 39);
    const imageOption = { offset: new window.kakao.maps.Point(27, 45) };

    const activeMarkerImg = new window.kakao.maps.MarkerImage(
      "/activeMarker.svg",
      imageSize,
      imageOption
    );

    const marker = new window.kakao.maps.Marker({
      image: activeMarkerImg,
      zIndex: 4,
    });

    if (latitude && longitude) {
      const markerPosition = new window.kakao.maps.LatLng(latitude, longitude);
      marker.setPosition(markerPosition);
    }
    marker.setMap(newMap);

    window.kakao.maps.event.addListener(
      newMap,
      "click",
      async (mouseEvent: KaKaoMapMouseEvent) => {
        const latlng = mouseEvent.latLng;

        const addr = (await getAddress(
          latlng.getLat(),
          latlng.getLng()
        )) as AddressInfo;

        setAddress(addr.address_name);

        marker.setPosition(latlng);
        marker.setVisible(true);

        setFormPosition(latlng.getLat(), latlng.getLng());
      }
    );
  }, [loaded]);

  return (
    <div
      className="relative"
      onMouseEnter={() => {
        if (isClickMap) return;
        setMapOver(true);
      }}
      onMouseLeave={() => {
        setMapOver(false);
        setIsClickMap(false);
      }}
    >
      {mapOver && (
        <div
          className="absolute top-0 left-0 w-full h-72 z-50 bg-black-tp-dark flex items-center justify-center text-2xl cursor-pointer"
          onClick={() => {
            setDraggable(true);
            setIsClickMap(true);
            setMapOver(false);
          }}
        >
          위치 선택하기
        </div>
      )}

      <div
        id="map"
        ref={mapRef}
        className={`w-full h-72 rounded-md overflow-hidden`}
      />
      <p className="text-sm mt-2">{address}</p>
    </div>
  );
};

export default MiniMap;
