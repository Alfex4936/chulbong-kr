"use client";

import useMapControl from "@/hooks/common/useMapControl";
import useMapStore from "@/store/useMapStore";
import { useSearchParams } from "next/navigation";
import { useEffect } from "react";

const isInSouthKorea = (lat: number, lng: number) => {
  const minLongitude = 124.6;
  const maxLongitude = 131.87;
  const minLatitude = 33.1;
  const maxLatitude = 38.45;

  return (
    lng >= minLongitude &&
    lng <= maxLongitude &&
    lat >= minLatitude &&
    lat <= maxLatitude
  );
};

const LocationQuery = () => {
  const { map } = useMapStore();
  const { moveLocation } = useMapControl();

  const searchParams = useSearchParams();
  const lat = searchParams.get("lat");
  const lng = searchParams.get("lng");

  useEffect(() => {
    if (!lat || !lng || !map || !isInSouthKorea(Number(lat), Number(lng))) {
      return;
    }
    moveLocation({
      lat: Number(lat),
      lng: Number(lng),
    });
  }, [map]);

  return null;
};

export default LocationQuery;
