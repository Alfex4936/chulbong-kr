"use client";

import useMiniMapStatusStore from "@/store/useMiniMapStatusStore";
import useMobileMapOpenStore from "@/store/useMobileMapOpenStore";
import useRoadviewStatusStore from "@/store/useRoadviewStatusStore";
import { useEffect, useState } from "react";
import Map from "./Map";
import MapLoading from "./MapLoading";
import Roadview from "./Roadview";

const MapWrapper = () => {
  const { isOpen } = useMobileMapOpenStore();
  const { isOpen: isRoadview } = useRoadviewStatusStore();
  const { setLoad } = useMiniMapStatusStore();

  const [loaded, setLoaded] = useState(false);

  useEffect(() => {
    const scriptEl = document.createElement("script");
    scriptEl.src = `//dapi.kakao.com/v2/maps/sdk.js?appkey=${process.env.NEXT_PUBLIC_APP_KEY}&libraries=clusterer,services&autoload=false`;
    scriptEl.async = true;

    document.head.appendChild(scriptEl);

    const handleLoadMap = () => {
      window.kakao.maps.load(() => {
        setLoaded(true);
        setLoad();
      });
    };

    scriptEl.addEventListener("load", handleLoadMap);

    return () => scriptEl.removeEventListener("load", handleLoadMap);
  }, []);

  return (
    <div
      className={`w-full h-dvh mo:absolute ${isOpen ? "mo:z-10" : "mo:-z-10"}`}
    >
      {!loaded ? (
        <div className="w-full h-dvh">
          <MapLoading />
        </div>
      ) : (
        <Map />
      )}

      {loaded && isRoadview && <Roadview />}
    </div>
  );
};

export default MapWrapper;
