import type { KakaoMap } from "@/types/KakaoMap.types";
import { MutableRefObject, useEffect, useState } from "react";
import useMapPositionStore from "../store/useMapPositionStore";

const useMap = (ref: MutableRefObject<HTMLDivElement | null>) => {
  const positionState = useMapPositionStore();

  const [map, setMap] = useState<KakaoMap | null>(null);

  useEffect(() => {
    const setNewPosition = (lat: number, lng: number) => {
      const options = {
        center: new window.kakao.maps.LatLng(lat, lng),
        level: positionState.level,
        maxLevel: 12,
      };

      const map = new window.kakao.maps.Map(ref.current, options);

      window.kakao.maps.event.addListener(map, "dragend", () => {
        const latlng = map.getCenter();

        positionState.setPosition(latlng.getLat(), latlng.getLng());
      });

      window.kakao.maps.event.addListener(map, "zoom_changed", () => {
        const level = map.getLevel();

        positionState.setLevel(level);
      });

      setMap(map);
    };

    if (navigator.geolocation) {
      navigator.geolocation.getCurrentPosition(
        (position) => {
          if (positionState.lat === 0 || positionState.lng === 0) {
            setNewPosition(position.coords.latitude, position.coords.longitude);
          } else {
            setNewPosition(positionState.lat, positionState.lng);
          }
        },
        (error) => {
          console.error(error);
        }
      );
    } else {
      if (positionState.lat === 0 || positionState.lng === 0) {
        setNewPosition(37.566535, 126.9779692);
      } else {
        setNewPosition(positionState.lat, positionState.lng);
      }
    }
  }, []);

  return map;
};

export default useMap;
