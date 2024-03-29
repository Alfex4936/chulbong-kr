import { Marker } from "@/types/Marker.types";
import instance from "../instance";

export type MarkerRes = Pick<Marker, "markerId" | "latitude" | "longitude">;

const getAllMarker = async (): Promise<MarkerRes[]> => {
  try {
    const res = await instance.get(`/api/v1/markers`);

    return res.data;
  } catch (error) {
    throw error;
  }
};

export default getAllMarker;
