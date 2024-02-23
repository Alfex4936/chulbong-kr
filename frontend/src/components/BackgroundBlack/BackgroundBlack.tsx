import * as Styled from "./BackgroundBlack.style";

interface Props {
  children: React.ReactNode;
}

const BackgroundBlack = ({ children }: Props) => {
  return (
    <Styled.BlackContainer>
      <Styled.ChildContainer>{children}</Styled.ChildContainer>
    </Styled.BlackContainer>
  );
};

export default BackgroundBlack;
