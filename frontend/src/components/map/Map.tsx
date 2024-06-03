"use client";

import { MAP_LAT_DIF } from "@/constants";
import useMyGps from "@/hooks/common/useMyGps";
import useAllMarkerData from "@/hooks/query/marker/useAllMarkerData";
import useBodyToggleStore from "@/store/useBodyToggleStore";
import useMapStatusStore from "@/store/useMapStatusStore";
import useMapStore from "@/store/useMapStore";
import useMobileMapOpenStore from "@/store/useMobileMapOpenStore";
import useSelectedMarkerStore from "@/store/useSelectedMarkerStore";
import { usePathname, useRouter } from "next/navigation";
import { useEffect, useRef, useState } from "react";
import { createRoot } from "react-dom/client";
import { MdOutlineGpsFixed } from "react-icons/md";
import MinusIcon from "../icons/MinusIcon";
import PlusIcon from "../icons/PlusIcon";
import MarkerOverlay from "../layout/MarkerOverlay";
import RQProvider from "../provider/RQProvider";
import { Separator } from "../ui/separator";
import MapButtons from "./MapButtons";
import MapLoading from "./MapLoading";
import MapSearch from "./MapSearch";

const Map = () => {
  const pathname = usePathname();
  const router = useRouter();

  const path1 = pathname.split("/")[1];

  const { marker: selectedMarker } = useSelectedMarkerStore();
  const { centerMapOnCurrentPosition } = useMyGps();

  const { isOpen: isMobileMapOpen } = useMobileMapOpenStore();
  const { lat, lng, level, setLevel, setPosition } = useMapStatusStore();

  const {
    map,
    setMap,
    clusterer: sClusterer,
    setClusterer,
    setMarkers,
    setOverlay,
    markers: mapMarkers,
  } = useMapStore();
  const { isOpen } = useBodyToggleStore();

  const { data: markers } = useAllMarkerData();

  const [mapLoading, setMapLoading] = useState(true);

  const mapRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (!markers) return;

    const imageSize = new window.kakao.maps.Size(39, 39);
    const imageOption = { offset: new window.kakao.maps.Point(27, 45) };

    const activeMarkerImg = new window.kakao.maps.MarkerImage(
      "/activeMarker.svg",
      imageSize,
      imageOption
    );

    const options = {
      center: new window.kakao.maps.LatLng(lat, lng),
      level: level,
      maxLevel: 12,
    };

    const newMap = map || new window.kakao.maps.Map(mapRef.current, options);

    const handleDrag = () => {
      const latlng = newMap.getCenter();
      setPosition(latlng.getLat(), latlng.getLng());
    };

    const handleZoom = () => {
      const level = newMap.getLevel();
      setLevel(level);
    };

    window.kakao.maps.event.addListener(newMap, "dragend", handleDrag);
    window.kakao.maps.event.addListener(newMap, "zoom_changed", handleZoom);

    const clusterer =
      sClusterer ||
      new window.kakao.maps.MarkerClusterer({
        minClusterSize: 3,
        map: newMap,
        averageCenter: true,
        minLevel: 6,
        calculator: [10, 30, 50, 100],
        styles: [
          {
            width: "30px",
            height: "30px",
            background: "rgba(51, 204, 255, .8)",
            borderRadius: "15px",
            color: "#000",
            textAlign: "center",
            fontWeight: "bold",
            lineHeight: "31px",
            userSelect: "none",
          },
          {
            width: "40px",
            height: "40px",
            background: "rgba(255, 153, 0, .8)",
            borderRadius: "20px",
            color: "#000",
            textAlign: "center",
            fontWeight: "bold",
            lineHeight: "41px",
            userSelect: "none",
          },
          {
            width: "50px",
            height: "50px",
            background: "rgba(255, 51, 204, .8)",
            borderRadius: "25px",
            color: "#000",
            textAlign: "center",
            fontWeight: "bold",
            lineHeight: "51px",
            userSelect: "none",
          },
          {
            width: "60px",
            height: "60px",
            background: "rgba(98, 208, 111, .8)",
            borderRadius: "30px",
            color: "#000",
            textAlign: "center",
            fontWeight: "bold",
            lineHeight: "61px",
            userSelect: "none",
          },
          {
            width: "60px",
            height: "60px",
            background: "rgba(227, 103, 72, .8)",
            borderRadius: "30px",
            color: "#000",
            textAlign: "center",
            fontWeight: "bold",
            lineHeight: "61px",
            userSelect: "none",
          },
        ],
        zIndex: 8,
      });

    if (sClusterer) sClusterer.clear();

    const overlayDiv = document.createElement("div");
    overlayDiv.classList.add("overlay_1");
    const root = createRoot(overlayDiv);

    const newMarkers = markers?.map((marker) => {
      const newMarker = new window.kakao.maps.Marker({
        position: new window.kakao.maps.LatLng(
          marker.latitude,
          marker.longitude
        ),
        image: activeMarkerImg,
        title: marker.markerId,
        zIndex: 4,
      });

      window.kakao.maps.event.addListener(newMarker, "click", async () => {
        const moveLatLon = new window.kakao.maps.LatLng(
          (marker.latitude as number) + MAP_LAT_DIF,
          marker.longitude
        );

        newMap?.panTo(moveLatLon);
        setPosition((marker.latitude as number) + MAP_LAT_DIF, marker.longitude);

        if (document.getElementsByClassName("overlay_1")[0]) {
          document.getElementsByClassName("overlay_1")[0].remove();
        }

        const latlng = new window.kakao.maps.LatLng(
          marker.latitude,
          marker.longitude
        );

        const overlay = new window.kakao.maps.CustomOverlay({
          position: latlng,
          content: overlayDiv,
          zIndex: 11,
        });

        setOverlay(overlay);
        overlay.setMap(newMap);

        const closeOverlay = () => {
          overlay.setMap(null);
        };

        const goDetail = () => {
          router.push(`/pullup/${marker.markerId}`);
        };
        const goReport = () => {
          router.push(`/pullup/${marker.markerId}/report`);
        };

        root.render(
          <RQProvider>
            <MarkerOverlay
              markerId={marker.markerId}
              closeOverlay={closeOverlay}
              goDetail={goDetail}
              goReport={goReport}
              lat={marker.latitude}
              lng={marker.longitude}
            />
          </RQProvider>
        );
      });

      return newMarker;
    });

    clusterer.addMarkers(newMarkers);

    setMapLoading(false);
    setMap(newMap);
    setClusterer(clusterer);
    setMarkers(newMarkers);
    return () => {
      window.kakao.maps.event.removeListener(newMap, "dragend", handleDrag);
      window.kakao.maps.event.removeListener(
        newMap,
        "zoom_changed",
        handleZoom
      );
    };
  }, [markers]);

  useEffect(() => {
    if (!map || !mapMarkers) return;

    let resizeTime: NodeJS.Timeout;

    if (path1 === "pullup") {
      if (!selectedMarker) return;

      const filterClickMarker = () => {
        const imageSize = new window.kakao.maps.Size(39, 39);
        const imageOption = { offset: new window.kakao.maps.Point(27, 45) };
        const selectedMarkerImg = new window.kakao.maps.MarkerImage(
          "/selectedMarker.svg",
          imageSize,
          imageOption
        );
        const activeMarkerImg = new window.kakao.maps.MarkerImage(
          "/activeMarker.svg",
          imageSize,
          imageOption
        );
        mapMarkers.forEach((marker) => {
          if (Number(marker.getTitle()) === selectedMarker.markerId) {
            marker.setImage(selectedMarkerImg);
          } else {
            marker.setImage(activeMarkerImg);
          }
        });
      };

      const moveLatLon = new window.kakao.maps.LatLng(
        (lat as number) + MAP_LAT_DIF,
        lng
      );
      map.relayout();

      resizeTime = setTimeout(() => {
        map.setCenter(moveLatLon);
        map.relayout();
        filterClickMarker();
      }, 200);
    } else {
      // const moveLatLon = new window.kakao.maps.LatLng(
      //   (lat as number) + 0.003,
      //   lng
      // );
      // map.relayout();
      // resizeTime = setTimeout(() => {
      //   map.setCenter(moveLatLon);
      //   map.relayout();
      // }, 200);
    }
    return () => clearTimeout(resizeTime);
  }, [mapLoading, map, isMobileMapOpen, path1, selectedMarker]);

  useEffect(() => {
    if (!map) return;

    const resizeTime = setTimeout(() => {
      map.relayout();
    }, 200);

    return () => clearTimeout(resizeTime);
  }, [isOpen]);

  useEffect(() => {
    if (!map) return;
    const moveLatLon = new window.kakao.maps.LatLng(
      (lat as number) + MAP_LAT_DIF,
      lng
    );
    map.relayout();
    map.setCenter(moveLatLon);
  }, [map]);

  const zoomIn = () => {
    const level = map?.getLevel();

    map?.setLevel((level as number) - 1);
  };

  const zoomOut = () => {
    const level = map?.getLevel();

    map?.setLevel((level as number) + 1);
  };

  return (
    <div className="relative w-full h-dvh">
      {mapLoading && <MapLoading />}
      <div
        id="map"
        ref={mapRef}
        className={`absolute top-0 left-0 w-full h-full ${
          mapLoading ? "hidden" : "block"
        }`}
      >
        <MapSearch />
        <MapButtons
          icon={<MdOutlineGpsFixed />}
          className="top-16 right-2"
          onClick={centerMapOnCurrentPosition}
          tooltipText="내 위치"
        />
        <div className="absolute top-[100px] right-2 bg-black-light-2 flex flex-col z-50 rounded-sm">
          <MapButtons
            icon={<PlusIcon size={20} />}
            className="static"
            onClick={zoomIn}
            tooltipText="확대"
          />
          <div className="px-1">
            <Separator className="bg-grey-dark-1" />
          </div>
          <MapButtons
            icon={<MinusIcon size={20} />}
            className="static"
            onClick={zoomOut}
            tooltipText="축소"
          />
        </div>
      </div>
    </div>
  );
};

export default Map;
