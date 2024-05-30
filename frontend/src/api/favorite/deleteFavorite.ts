import instance from "../instance";

const deleteFavorite = async (id: number) => {
  const res = await instance.delete(`/api/v1/markers/${id}/favorites`);

  return res.data;
};

export default deleteFavorite;
