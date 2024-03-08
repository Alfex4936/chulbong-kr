import { Marker } from "@/types/Marker.types";
import axios from "axios";

const getMarker = async (id: number): Promise<Marker> => {
  try {
    const res = await axios.get(`/api/v1/markers/${id}`);

    return res.data;
  } catch (error) {
    throw error;
  }
};

export default getMarker;
