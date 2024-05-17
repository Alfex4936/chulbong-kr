import instance from "../instance";

const deleteComment = async (id: number) => {
  const res = await instance.delete(`/api/v1/comments/${id}`);

  return res.data;
};

export default deleteComment;
