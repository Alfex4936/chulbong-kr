import { Marker } from "@/types/Marker.types";
import axios from "axios";

type MarkerRes = Pick<Marker, "markerId" | "latitude" | "longitude">;

const getAllMarker = async (): Promise<MarkerRes[]> => {
  try {
    const res = await axios.get(`/api/v1/markers`, {
      headers: { "Accept-Encoding": "gzip, deflate, br" },
    });

    return res.data;
  } catch (error) {
    throw error;
  }
};

export default getAllMarker;
