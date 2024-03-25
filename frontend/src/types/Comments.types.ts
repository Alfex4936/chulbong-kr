export interface Comment {
  commentId: number;
  markerId: number;
  userId: number;
  commentText: string;
  postedAt: Date;
  updatedAt: Date;
  username: string;
}
