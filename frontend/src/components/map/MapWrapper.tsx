"use client";

import Script from "next/script";
import { useState } from "react";
import Map from "./Map";
import MapLoading from "./MapLoading";

const MapWrapper = () => {
  const [loaded, setLoaded] = useState(false);

  return (
    <>
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
    </>
  );
};

export default MapWrapper;
