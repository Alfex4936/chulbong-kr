import { Marker } from "@/types/Marker.types";
import axios from "axios";

// : Promise<Marker[]>
interface MyMarkerRes {
  currentPage: number;
  markers: Marker[];
  totalMarkers: number;
  totalPages: number;
}

const getMyMarker = async ({
  pageParam,
}: {
  pageParam: number;
}): Promise<MyMarkerRes> => {
  try {
    const res = await axios.get(`/api/v1/markers/my?page=${pageParam}`);

    return res.data;
  } catch (error) {
    throw error;
  }
};

export default getMyMarker;
