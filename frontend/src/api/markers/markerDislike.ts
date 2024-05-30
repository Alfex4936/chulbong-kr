import instance from "../instance";

const markerDislike = async (
  markerId: number
): Promise<{ disliked: boolean }> => {
  const res = await instance.post(`/api/v1/markers/${markerId}/dislike`);

  return res.data;
};

export default markerDislike;
