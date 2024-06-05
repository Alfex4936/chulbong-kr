import { Marker } from "@/types/Marker.types";
import instance from "../instance";

export type MarkerRes = Pick<
  Marker,
  "markerId" | "latitude" | "longitude" | "address"
>;

const getAllMarker = async (): Promise<MarkerRes[]> => {
  const res = await instance.get(`/api/v1/markers`);

  return res.data;
};

export default getAllMarker;
