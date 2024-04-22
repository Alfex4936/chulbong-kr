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
        src="//dapi.kakao.com/v2/maps/sdk.js?appkey=dfdebaf84d7dda475fb8448c7d43c528&libraries=clusterer,services&autoload=false"
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
