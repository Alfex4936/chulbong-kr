import axios from "axios";

interface setMarkerReq {
  photos: File;
  latitude: number;
  longitude: number;
  description: string;
}

const setNewMarker = async (multipart: setMarkerReq) => {
  const formData = new FormData();
  const token = JSON.parse(localStorage.getItem("user") as string).state.user
    .token;

  formData.append("photos", multipart.photos);
  formData.append("latitude", multipart.latitude.toString());
  formData.append("longitude", multipart.longitude.toString());
  formData.append("description", multipart.description);

  try {
    const res = await axios.post(`/api/v1/markers/new`, formData, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });

    return res;
  } catch (error) {
    throw error;
  }
};

export default setNewMarker;
