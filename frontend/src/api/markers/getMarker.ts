import type { Marker } from "../../types/Marker.types";
import instance from "../instance";

const getMarker = async (id: number): Promise<Marker> => {
  const res = await instance.get(`/api/v1/markers/${id}/details`);

  return res.data;
};

export default getMarker;
