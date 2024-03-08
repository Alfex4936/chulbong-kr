import * as Styled from "./MarkerReviewSkeleton.style";

const MarkerReviewSkeleton = () => {
  return (
    <>
      <Styled.ListSkeleton>
        <div />
        <div style={{ flexGrow: "1" }} />
        <div />
      </Styled.ListSkeleton>
      <Styled.ListSkeleton>
        <div />
        <div style={{ flexGrow: "1" }} />
        <div />
      </Styled.ListSkeleton>
      <Styled.ListSkeleton>
        <div />
        <div style={{ flexGrow: "1" }} />
        <div />
      </Styled.ListSkeleton>
      <Styled.ListSkeleton>
        <div />
        <div style={{ flexGrow: "1" }} />
        <div />
      </Styled.ListSkeleton>
    </>
  );
};

export default MarkerReviewSkeleton;
