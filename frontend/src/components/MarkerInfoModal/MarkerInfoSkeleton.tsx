import * as Styled from "./MarkerInfoSkeleton.style";

const MarkerInfoSkeleton = () => {
  return (
    <div>
      <Styled.imageWrap>
        <Styled.SkeletonImage />
      </Styled.imageWrap>
      <Styled.SkeletonButtons />
    </div>
  );
};

export default MarkerInfoSkeleton;
