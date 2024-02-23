import axios from "axios";

const getAllMarker = async () => {
  try {
    const res = await axios.get(`/api/v1/markers`);

    return res;
  } catch (error) {
    throw error;
  }
};

export default getAllMarker;
