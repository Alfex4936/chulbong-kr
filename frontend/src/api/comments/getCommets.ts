import type { Comment } from "../../types/Comments.types";
import instance from "../instance";

interface Props {
  id: number;
  pageParam: number;
}

interface CommentsRes {
  currentPage: number;
  comments: Comment[];
  totalComments: number;
  totalPages: number;
}

const getCommets = async ({ id, pageParam }: Props): Promise<CommentsRes> => {
  try {
    const res = await instance.get(
      `/api/v1/comments/${id}/comments?page=${pageParam}`
    );

    return res.data;
  } catch (error) {
    throw error;
  }
};

export default getCommets;
