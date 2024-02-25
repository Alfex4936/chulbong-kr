import { Marker } from "@/types/Marker.types";
import axios from "axios";

const getAllMarker = async (): Promise<Marker[]> => {
  try {
    const res = await axios.get(`/api/v1/markers`);

    return res.data;
  } catch (error) {
    throw error;
  }
};

export default getAllMarker;
