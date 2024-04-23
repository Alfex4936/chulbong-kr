import instance from "../instance";

const setFavorite = async (id: number) => {
  const res = await instance.post(`/api/v1/markers/${id}/favorites`);

  return res.data;
};

export default setFavorite;
