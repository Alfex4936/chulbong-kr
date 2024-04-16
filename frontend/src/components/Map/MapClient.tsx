"use client";

import useAllMarkerData from "@/hooks/query/useAllMarkerData";
import useBodyToggleStore from "@/store/useBodyToggleStore";
import useMapStore from "@/store/useMapStore";
import { usePathname, useRouter } from "next/navigation";
import { useEffect, useRef, useState } from "react";
import MapLoading from "./MapLoading";
import useMapStatusStore from "@/store/useMapStatusStore";

const MapClient = () => {
  const { lat, lng, level, setLevel, setPosition } = useMapStatusStore();

  const { map, setMap, setClusterer } = useMapStore();
  const { isOpen } = useBodyToggleStore();

  const [mapLoading, setMapLoading] = useState(true);

  const { data: markers } = useAllMarkerData();

  const mapRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (!markers) return;

    const nlat = lat || 37.566535;
    const nlng = lng || 126.9779692;
    const nlevel = level || 3;

    console.log(lat, lng, level);
    const options = {
      center: new window.kakao.maps.LatLng(lat, lng),
      level: level,
      maxLevel: 12,
    };

    const newMap = new window.kakao.maps.Map(mapRef.current, options);

    const handleDrag = () => {
      // const level = map.getLevel();
      const latlng = newMap.getCenter();
      // const query = `lat=${latlng.getLat()}&lng=${latlng.getLng()}&lv=${level}`;
      // router.push(`${pathname}?${query}`);
      setPosition(latlng.getLat(), latlng.getLng());
    };

    const handleZoom = () => {
      const level = newMap.getLevel();
      // const latlng = map.getCenter();
      // const query = `lat=${latlng.getLat()}&lng=${latlng.getLng()}&lv=${level}`;

      // router.push(`${pathname}?${query}`);
      setLevel(level);
    };

    window.kakao.maps.event.addListener(newMap, "dragend", handleDrag);
    window.kakao.maps.event.addListener(newMap, "zoom_changed", handleZoom);

    const imageSize = new window.kakao.maps.Size(39, 39);
    const imageOption = { offset: new window.kakao.maps.Point(27, 45) };

    const activeMarkerImg = new window.kakao.maps.MarkerImage(
      "/activeMarker.svg",
      imageSize,
      imageOption
    );

    const clusterer = new window.kakao.maps.MarkerClusterer({
      map: map,
      averageCenter: true,
      minLevel: 6,
    });

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

      return newMarker;
    });

    clusterer.addMarkers(newMarkers);

    setMapLoading(false);
    setMap(newMap);
    setClusterer(clusterer);

    return () => {
      window.kakao.maps.event.removeListener(newMap, "dragend", handleDrag);
      window.kakao.maps.event.removeListener(
        newMap,
        "zoom_changed",
        handleZoom
      );
    };
  }, [markers]);

  // useEffect(() => {
  //   if (!map) return;
  // }, []);

  useEffect(() => {
    if (!map) return;

    map.relayout();

    const resizeTime = setTimeout(() => map.relayout(), 200);

    return () => clearTimeout(resizeTime);
  }, [isOpen, mapLoading, map]);

  return (
    <div className="relative w-full mo:hidden">
      {mapLoading && <MapLoading />}
      <div
        id="map"
        ref={mapRef}
        className={`absolute top-0 left-0 w-full h-full ${
          mapLoading ? "hidden" : "block"
        }`}
      />
    </div>
  );
};

export default MapClient;
