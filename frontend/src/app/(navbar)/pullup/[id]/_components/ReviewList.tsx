"use client";

import GrowBox from "@/components/atom/GrowBox";
import DeleteIcon from "@/components/icons/DeleteIcon";
import useGetComments from "@/hooks/query/comments/useCommentsData";
import { useEffect, useRef } from "react";
import formatDate from "@/utils/formatDate";
import { Skeleton } from "@/components/ui/skeleton";
// TODO: 유저 정보 삭제 연동
// TODO: 리뷰 리스트 스타일 변경

type Props = {
  markerId: number;
};

const ReviewList = ({ markerId }: Props) => {
  const {
    data: review,
    fetchNextPage,
    hasNextPage,
    isLoading,
    isError,
    isFetching,
  } = useGetComments(markerId);

  const boxRef = useRef(null);

  useEffect(() => {
    const currentRef = boxRef.current;

    const observer = new IntersectionObserver(
      (entries) => {
        const entry = entries[0];
        if (entry.isIntersecting) {
          if (!isFetching && hasNextPage) {
            fetchNextPage();
          }
        }
      },
      { threshold: 0.8 }
    );

    if (currentRef) {
      observer.observe(currentRef);
    }

    return () => {
      if (currentRef) {
        observer.unobserve(currentRef);
      }
    };
  }, [isFetching, hasNextPage, fetchNextPage]);

  if (isError) return <div>잠시 후 다시 시도해 주세요!</div>;

  if (isLoading) {
    return <Skeleton className="w-full h-12 rounded-sm bg-black-light-2" />;
  }

  return (
    <div>
      {review?.pages[0].comments.length === 0 && (
        <div>등록된 리뷰가 없습니다.</div>
      )}

      {review?.pages.map((page, index) => {
        return (
          <div key={index}>
            {page.comments.map((comment) => {
              return (
                <div className="flex items-center p-3" key={comment.commentId}>
                  <div className="w-2/3">
                    <div className="truncate text-xl mr-2 overflow-hidden whitespace-nowrap hover:whitespace-normal hover:overflow-visible hover:break-words">
                      {comment.commentText}
                    </div>
                    <div className="truncate text-[10px] text-grey-dark-1">
                      {formatDate(comment.postedAt)}
                    </div>
                  </div>
                  <GrowBox />
                  <span className="text-xs">{comment.username}</span>
                  <button className="ml-2">
                    <DeleteIcon />
                  </button>
                </div>
              );
            })}
          </div>
        );
      })}

      {hasNextPage && <div ref={boxRef} />}
    </div>
  );
};

export default ReviewList;
