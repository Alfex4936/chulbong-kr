import axios from "axios";

const deleteComment = async (id: number) => {
  try {
    const res = await axios.delete(`/api/v1/comments/${id}`);

    return res.data;
  } catch (error) {
    throw error;
  }
};

export default deleteComment;
