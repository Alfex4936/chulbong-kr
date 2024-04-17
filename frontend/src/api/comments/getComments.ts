import type { Comment } from "../../types/Comments.types";
import instance from "../instance";

interface Props {
  id: number;
  pageParam: number;
}

export interface CommentsRes {
  currentPage: number;
  comments: Comment[];
  totalComments: number;
  totalPages: number;
}

const getComments = async ({ id, pageParam }: Props): Promise<CommentsRes> => {
  const res = await instance.get(
    `/api/v1/comments/${id}/comments?page=${pageParam}`
  );

  return res.data;
};

export default getComments;
