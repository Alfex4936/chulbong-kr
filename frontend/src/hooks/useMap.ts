import type { KakaoMap } from "@/types/KakaoMap.types";
import { MutableRefObject, useEffect, useState } from "react";

const useMap = (ref: MutableRefObject<HTMLDivElement | null>) => {
  const [map, setMap] = useState<KakaoMap | null>(null);

  useEffect(() => {
    var options = {
      center: new window.kakao.maps.LatLng(37.566535, 126.9779692),
      level: 3,
    };

    var map = new window.kakao.maps.Map(ref.current, options);

    setMap(map);
  }, []);

  return map;
};

export default useMap;
