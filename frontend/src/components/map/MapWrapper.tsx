"use client";

import useMiniMapStatusStore from "@/store/useMiniMapStatusStore";
import useMobileMapOpenStore from "@/store/useMobileMapOpenStore";
import useRoadviewStatusStore from "@/store/useRoadviewStatusStore";
import Script from "next/script";
import { useState } from "react";
import Map from "./Map";
import MapLoading from "./MapLoading";
import Roadview from "./Roadview";

const MapWrapper = () => {
  const { isOpen } = useMobileMapOpenStore();
  const { isOpen: isRoadview } = useRoadviewStatusStore();
  const { setLoad } = useMiniMapStatusStore();

  const [loaded, setLoaded] = useState(false);

  return (
    <div
      className={`w-full h-screen mo:absolute ${
        isOpen ? "mo:z-10" : "mo:-z-10"
      }`}
    >
      <Script
        src={`//dapi.kakao.com/v2/maps/sdk.js?appkey=${process.env.NEXT_PUBLIC_APP_KEY}&libraries=clusterer,services&autoload=false`}
        onLoad={() => {
          window.kakao.maps.load(() => {
            setLoaded(true);
            setLoad();
          });
        }}
      />

      {!loaded && (
        <div className="w-full h-screen">
          <MapLoading />
        </div>
      )}
      {loaded && <Map />}

      {loaded && isRoadview && <Roadview />}
    </div>
  );
};

export default MapWrapper;
