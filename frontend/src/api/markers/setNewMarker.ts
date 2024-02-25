import { Marker } from "@/types/Marker.types";
import axios from "axios";

export interface SetMarkerReq {
  photos: File;
  latitude: number;
  longitude: number;
  description: string;
}

export interface SetMarkerRes
  extends Omit<Marker, "photos" | "createdAt" | "updatedAt"> {
  photoUrls?: string[];
}

const setNewMarker = async (multipart: SetMarkerReq): Promise<SetMarkerRes> => {
  const formData = new FormData();

  formData.append("photos", multipart.photos);
  formData.append("latitude", multipart.latitude.toString());
  formData.append("longitude", multipart.longitude.toString());
  formData.append("description", multipart.description);

  try {
    const res = await axios.post(`/api/v1/markers/new`, formData, {
      withCredentials: true,
    });

    return res.data;
  } catch (error) {
    throw error;
  }
};

export default setNewMarker;
