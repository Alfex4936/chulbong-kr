import instance from "../instance";

const createComment = async (body: {
  markerId: number;
  commentText: string;
}) => {
  try {
    const res = await instance.post(`/api/v1/comments`, body);

    return res.data;
  } catch (error) {
    throw error;
  }
};

export default createComment;
