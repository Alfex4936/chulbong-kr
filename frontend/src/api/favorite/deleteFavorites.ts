import axios from "axios";

const deleteFavorites = async (id: number) => {
  try {
    const res = await axios.delete(`/api/v1/markers/${id}/favorites`);

    return res.data;
  } catch (error) {
    throw error;
  }
};

export default deleteFavorites;
