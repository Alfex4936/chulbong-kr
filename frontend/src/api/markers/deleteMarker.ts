import axios from "axios";

const deleteMarker = async (id: number) => {
  try {
    const res = await axios.delete(`/api/v1/markers/${id}`, {
      withCredentials: true,
    });

    return res;
  } catch (error) {
    throw error;
  }
};

export default deleteMarker;
