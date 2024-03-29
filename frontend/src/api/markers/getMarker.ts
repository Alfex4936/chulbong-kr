import type { Marker } from "../../types/Marker.types";
import instance from "../instance";

const getMarker = async (id: number): Promise<Marker> => {
  try {
    const res = await instance.get(`/api/v1/markers/${id}/details`);

    return res.data;
  } catch (error) {
    throw error;
  }
};

export default getMarker;
