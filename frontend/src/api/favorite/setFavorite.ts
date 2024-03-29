import instance from "../instance";

const setFavorite = async (id: number) => {
  try {
    const res = await instance.post(`/api/v1/markers/${id}/favorites`);

    return res.data;
  } catch (error) {
    throw error;
  }
};

export default setFavorite;
