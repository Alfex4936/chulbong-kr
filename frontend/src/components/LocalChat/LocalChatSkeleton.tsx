import * as Styled from "./LocalChatSkeleton.style";

const LocalChatSkeleton = () => {
  return (
    <Styled.Container>
      <Styled.SkeletonBox />
      <Styled.SkeletonInput />
    </Styled.Container>
  );
};

export default LocalChatSkeleton;
