import axios from "axios";

const getDislikeState = async (
  markerId: number
): Promise<{ disliked: boolean }> => {
  try {
    const res = await axios.get(`/api/v1/markers/${markerId}/dislike-status`);

    return res.data;
  } catch (error) {
    throw error;
  }
};

export default getDislikeState;
