import type { Marker } from "@/types/Marker.types";
import instance from "../instance";

export interface SetMarkerReq {
  photos: File[];
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

  for (let i = 0; i < multipart.photos.length; i++) {
    formData.append("photos", multipart.photos[i]);
  }

  formData.append("latitude", multipart.latitude.toString());
  formData.append("longitude", multipart.longitude.toString());
  formData.append("description", multipart.description);

  const res = await instance.post(`/api/v1/markers/new`, formData, {
    withCredentials: true,
  });

  return res.data;
};

export default setNewMarker;
