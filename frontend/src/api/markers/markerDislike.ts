import instance from "../instance";

const markerDislike = async (
  markerId: number
): Promise<{ disliked: boolean }> => {
  try {
    const res = await instance.post(`/api/v1/markers/${markerId}/dislike`);

    return res.data;
  } catch (error) {
    throw error;
  }
};

export default markerDislike;
