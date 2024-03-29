import instance from "../instance";

const deleteComment = async (id: number) => {
  try {
    const res = await instance.delete(`/api/v1/comments/${id}`);

    return res.data;
  } catch (error) {
    throw error;
  }
};

export default deleteComment;
