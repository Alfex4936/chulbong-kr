import { Marker } from "@/types/Marker.types";
import instance from "../instance";

// ${process.env.NEXT_PUBLIC_BASE_URL}
export type MarkerRes = Pick<Marker, "markerId" | "latitude" | "longitude">;

const getAllMarker = async (): Promise<MarkerRes[]> => {
  const res = await instance.get(`/api/v1/markers`);

  return res.data;
};

export default getAllMarker;
