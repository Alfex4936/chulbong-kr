"use client";

import { MOBILE_WIDTH } from "@/constants";
import useMiniMapStatusStore from "@/store/useMiniMapStatusStore";
import useMobileMapOpenStore from "@/store/useMobileMapOpenStore";
import useRoadviewStatusStore from "@/store/useRoadviewStatusStore";
import Script from "next/script";
import { useEffect, useState } from "react";
import Map from "./Map";
import MapLoading from "./MapLoading";
import Roadview from "./Roadview";

const MapWrapper = () => {
  const { isOpen, open, close } = useMobileMapOpenStore();
  const { isOpen: isRoadview } = useRoadviewStatusStore();
  const { setLoad } = useMiniMapStatusStore();

  const [loaded, setLoaded] = useState(false);

  useEffect(() => {
    const handleResize = () => {
      if (window.innerWidth > MOBILE_WIDTH) {
        open();
      } else {
        close();
      }
    };

    handleResize();

    window.addEventListener("resize", handleResize);

    return () => {
      window.removeEventListener("resize", handleResize);
    };
  }, []);

  return (
    <div
      className={`w-full h-dvh mo:absolute ${isOpen ? "mo:z-10" : "mo:-z-10"}`}
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
        <div className="w-full h-dvh">
          <MapLoading />
        </div>
      )}
      {loaded && <Map />}

      {loaded && isRoadview && <Roadview />}

      {!isOpen && (
        <div className="absolute top-0 left-0 w-dvw h-dvh bg-black z-30" />
      )}
    </div>
  );
};

export default MapWrapper;
