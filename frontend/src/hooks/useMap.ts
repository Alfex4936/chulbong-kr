import type { KakaoMap } from "@/types/KakaoMap.types";
import { MutableRefObject, useEffect, useState } from "react";
import useMapPositionStore from "../store/useMapPositionStore";

const useMap = (ref: MutableRefObject<HTMLDivElement | null>) => {
  const positionState = useMapPositionStore();

  const [map, setMap] = useState<KakaoMap | null>(null);

  useEffect(() => {
    const options = {
      center: new window.kakao.maps.LatLng(
        positionState.lat,
        positionState.lng
      ),
      level: positionState.level,
      maxLevel: 12,
    };

    const map = new window.kakao.maps.Map(ref.current, options);

    const handleDrag = () => {
      const latlng = map.getCenter();

      positionState.setPosition(latlng.getLat(), latlng.getLng());
    };

    const handleZoom = () => {
      const level = map.getLevel();

      positionState.setLevel(level);
    };

    window.kakao.maps.event.addListener(map, "dragend", handleDrag);

    window.kakao.maps.event.addListener(map, "zoom_changed", handleZoom);

    setMap(map);

    positionState.setPosition(positionState.lat, positionState.lng);

    return () => {
      window.kakao.maps.event.removeListener(map, "dragend", handleDrag);
      window.kakao.maps.event.removeListener(map, "zoom_changed", handleZoom);
    };
  }, []);

  return map;
};

export default useMap;
