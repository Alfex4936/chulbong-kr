import ArrowBackIcon from "@mui/icons-material/ArrowBack";
import DeleteOutlineIcon from "@mui/icons-material/DeleteOutline";
import ReplyIcon from "@mui/icons-material/Reply";
import CircularProgress from "@mui/material/CircularProgress";
import IconButton from "@mui/material/IconButton";
import Tooltip from "@mui/material/Tooltip";
import { isAxiosError } from "axios";
import { useEffect, useRef, useState } from "react";
import useCreateComment from "../../hooks/mutation/comments/useCreateComment";
import useDeleteComment from "../../hooks/mutation/comments/useDeleteComment";
import useGetComments from "../../hooks/query/comments/useGetComments";
import useGetMyInfo from "../../hooks/query/user/useGetMyInfo";
import useInput from "../../hooks/useInput";
import * as Styled from "./MarkerReview.style";
import MarkerReviewSkeleton from "./MarkerReviewSkeleton";

interface Props {
  markerId: number;
  setIsReview: React.Dispatch<React.SetStateAction<boolean>>;
}

const MarkerReview = ({ markerId, setIsReview }: Props) => {
  const commentValue = useInput("");

  const { data: myInfo } = useGetMyInfo();
  const { data, fetchNextPage, hasNextPage, isLoading, isError, isFetching } =
    useGetComments(markerId);

  const { mutateAsync, isPending } = useCreateComment({
    markerId,
    commentText: commentValue.value,
  });

  const { mutate: deleteComment, isPending: deleteLoading } =
    useDeleteComment(markerId);

  const [errorText, setErrorText] = useState("");

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

  const handleComment = async () => {
    if (commentValue.value === "") return;

    try {
      await mutateAsync();
      commentValue.reset();
    } catch (error) {
      if (isAxiosError(error)) {
        if (error.response?.status === 400) {
          setErrorText("비속어 사용을 조심해 주세요!!");
        } else {
          setErrorText("잠시 후 다시 시도해 주세요.");
        }
      }
    }
  };

  const handleDelete = (id: number) => {
    deleteComment(id);
  };

  const handleKeyPress = (event: React.KeyboardEvent<HTMLInputElement>) => {
    if (event.key === "Enter") {
      handleComment();
    }
  };

  if (isError) <div>잠시 후 다시 시도해 주세요!</div>;

  return (
    <div>
      <Tooltip title="이전" arrow disableInteractive>
        <IconButton
          onClick={() => {
            setIsReview(false);
          }}
          aria-label="delete"
          sx={{
            position: "absolute",
            top: "0",
            left: "0",
          }}
        >
          <ArrowBackIcon />
        </IconButton>
      </Tooltip>

      <Styled.ReviewWrap>
        {isLoading || deleteLoading || isPending ? (
          <MarkerReviewSkeleton />
        ) : (
          <>
            {data?.pages[0].comments.length === 0 && (
              <Styled.Container>
                <Styled.Wrapper>
                  <Styled.P>텅</Styled.P>

                  <Styled.DotWrap>
                    <Styled.Dot1></Styled.Dot1>
                    <Styled.Dot2></Styled.Dot2>
                    <Styled.Dot3></Styled.Dot3>
                  </Styled.DotWrap>
                </Styled.Wrapper>
                <Styled.Text>
                  <p>등록된 리뷰가 없습니다.</p>
                </Styled.Text>
              </Styled.Container>
            )}
            {data?.pages.map((page, index) => {
              return (
                <div key={index}>
                  {page.comments.map((comment, index) => {
                    return (
                      <Styled.ReviewItem key={index}>
                        <div>{comment.commentText}</div>
                        <div />
                        <div>
                          <div>
                            {comment.postedAt
                              .toString()
                              .split("T")[0]
                              .replace(/-/g, ".")}
                          </div>
                          <div>{comment.username}</div>
                        </div>
                        <div>
                          {myInfo?.userId === comment.userId && (
                            <Tooltip title="삭제 하기" arrow disableInteractive>
                              <IconButton
                                onClick={() => {
                                  handleDelete(comment.commentId);
                                }}
                                aria-label="delete"
                              >
                                {deleteLoading ? (
                                  <CircularProgress color="inherit" size={20} />
                                ) : (
                                  <DeleteOutlineIcon />
                                )}
                              </IconButton>
                            </Tooltip>
                          )}
                        </div>
                      </Styled.ReviewItem>
                    );
                  })}
                </div>
              );
            })}

            {hasNextPage && (
              <Styled.ListSkeleton ref={boxRef}>
                <div />
                <div style={{ flexGrow: "1" }} />
                <div />
              </Styled.ListSkeleton>
            )}
          </>
        )}
      </Styled.ReviewWrap>
      <Styled.InputWrap>
        <Styled.ReviewInput
          type="text"
          name="reveiw-content"
          maxLength={40}
          value={commentValue.value}
          onChange={(e) => {
            if (commentValue.value.length >= 40) {
              setErrorText("40자 이내로 작성해 주세요!");
            } else {
              setErrorText("");
            }
            commentValue.onChange(e);
          }}
          onKeyDown={handleKeyPress}
        />
        <Tooltip title="등록" arrow disableInteractive>
          <IconButton onClick={handleComment} aria-label="send">
            <ReplyIcon />
          </IconButton>
        </Tooltip>
      </Styled.InputWrap>
      <Styled.ErrorBox>{errorText}</Styled.ErrorBox>
    </div>
  );
};

export default MarkerReview;
