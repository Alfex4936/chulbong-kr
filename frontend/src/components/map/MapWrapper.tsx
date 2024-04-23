"use client";

import useMobileMapOpenStore from "@/store/useMobileMapOpenStore";
import Script from "next/script";
import { useState } from "react";
import Map from "./Map";
import MapLoading from "./MapLoading";

const MapWrapper = () => {
  const { isOpen } = useMobileMapOpenStore();

  const [loaded, setLoaded] = useState(false);

  return (
    <div
      className={`w-full h-screen mo:absolute mo:z-10 ${
        isOpen ? "mo:block" : "mo:hidden"
      }`}
    >
      <Script
        src={`//dapi.kakao.com/v2/maps/sdk.js?appkey=${process.env.NEXT_PUBLIC_APP_KEY}&libraries=clusterer,services&autoload=false`}
        onLoad={() => {
          window.kakao.maps.load(() => {
            setLoaded(true);
          });
        }}
      />

      {!loaded && (
        <div className="w-full h-screen">
          <MapLoading />
        </div>
      )}
      {loaded && <Map />}
    </div>
  );
};

export default MapWrapper;
