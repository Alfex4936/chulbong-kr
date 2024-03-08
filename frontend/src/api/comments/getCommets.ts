import axios from "axios";

const getCommets = async (id: number) => {
  try {
    const res = await axios.get(`/api/v1/comments/${id}/comments`);

    return res.data;
  } catch (error) {
    throw error;
  }
};

export default getCommets;
