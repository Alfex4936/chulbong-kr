import instance from "../instance";

const undoMarkerDislike = async (
  markerId: number
): Promise<{ disliked: boolean }> => {
  const res = await instance.delete(`/api/v1/markers/${markerId}/dislike`);

  return res.data;
};

export default undoMarkerDislike;
