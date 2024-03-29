import instance from "../instance";

const markerUnDislike = async (
  markerId: number
): Promise<{ disliked: boolean }> => {
  try {
    const res = await instance.delete(`/api/v1/markers/${markerId}/dislike`);

    return res.data;
  } catch (error) {
    throw error;
  }
};

export default markerUnDislike;
