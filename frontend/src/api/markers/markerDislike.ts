import axios from "axios";

const markerDislike = async (
  markerId: number
): Promise<{ disliked: boolean }> => {
  try {
    const res = await axios.post(`/api/v1/markers/${markerId}/dislike`);

    return res.data;
  } catch (error) {
    throw error;
  }
};

export default markerDislike;
