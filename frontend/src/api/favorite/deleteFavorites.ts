import instance from "../instance";

const deleteFavorites = async (id: number) => {
  try {
    const res = await instance.delete(`/api/v1/markers/${id}/favorites`);

    return res.data;
  } catch (error) {
    throw error;
  }
};

export default deleteFavorites;
