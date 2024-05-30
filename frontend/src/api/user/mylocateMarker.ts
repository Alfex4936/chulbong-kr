import type { Marker } from "../../types/Marker.types";
import instance from "../instance";

export interface MyMarkerRes {
  currentPage: number;
  markers: Marker[];
  totalMarkers: number;
  totalPages: number;
}

const mylocateMarker = async ({
  pageParam,
}: {
  pageParam: number;
}): Promise<MyMarkerRes> => {
  const res = await instance.get(
    `/api/v1/markers/my?page=${pageParam}&pageSize=10`
  );

  return res.data;
};

export default mylocateMarker;
