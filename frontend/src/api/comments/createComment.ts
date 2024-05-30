import instance from "../instance";

const createComment = async (body: {
  markerId: number;
  commentText: string;
}) => {
  const res = await instance.post(`/api/v1/comments`, body);

  return res.data;
};

export default createComment;
