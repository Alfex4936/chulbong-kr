import { useEffect, useState } from "react";

const useMap = (ref) => {
  const [map, setMap] = useState(null);

  useEffect(() => {
    var options = {
      center: new kakao.maps.LatLng(37.566535, 126.9779692),
      level: 3,
    };

    var map = new kakao.maps.Map(ref.current, options);

    setMap(map);
  }, []);

  return map;
};

export default useMap;
