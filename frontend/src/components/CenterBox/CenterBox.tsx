import * as Styled from "./CenterBox.style";

interface Props {
  children: React.ReactNode;
  bg?: "transparent" | "black";
}

const CenterBox = ({ children, bg = "transparent" }: Props) => {
  return (
    <Styled.BlackContainer bg={bg}>
      <Styled.ChildContainer>{children}</Styled.ChildContainer>
    </Styled.BlackContainer>
  );
};

export default CenterBox;
