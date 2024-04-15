"use client";

import useAllMarkerData from "@/hooks/query/useAllMarkerData";
import useBodyToggleStore from "@/store/useBodyToggleStore";
import useMapStore from "@/store/useMapStore";
import { useEffect, useRef, useState } from "react";
import MapLoading from "./MapLoading";

const MapClient = () => {
  const { map, setMap, setClusterer } = useMapStore();
  const { isOpen } = useBodyToggleStore();
  const [mapLoading, setMapLoading] = useState(true);

  const { data: markers } = useAllMarkerData();

  const mapRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (map) return;

    window.kakao.maps.load(() => {
      const options = {
        center: new window.kakao.maps.LatLng(33.450701, 126.570667),
        level: 3,
        maxLevel: 12,
      };

      const map = new window.kakao.maps.Map(mapRef.current, options);

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
      setMap(map);
      setClusterer(clusterer);
    });
  }, []);

  useEffect(() => {
    if (!map) return;

    map.relayout();
  }, [isOpen, mapLoading]);

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
