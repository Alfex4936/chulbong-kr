import ArrowBackIcon from "@mui/icons-material/ArrowBack";
import IconButton from "@mui/material/IconButton";
import Tooltip from "@mui/material/Tooltip";
import useGetComments from "../../hooks/query/comments/useGetComments";
import MarkerReviewSkeleton from "./MarkerReviewSkeleton";

interface Props {
  markerId: number;
  setIsReview: React.Dispatch<React.SetStateAction<boolean>>;
}

const MarkerReview = ({ markerId, setIsReview }: Props) => {
  const { data, isLoading, isError } = useGetComments(markerId);

  if (isLoading) return <MarkerReviewSkeleton />;
  if (isError) return <div>잠시 후 다시 시도해 주세요!</div>;

  console.log(data);

  return (
    <div>
      <Tooltip title="닫기" arrow disableInteractive>
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
      <div>리뷰</div>
      <div>리뷰</div>
    </div>
  );
};

export default MarkerReview;
