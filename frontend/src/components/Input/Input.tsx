import * as Styled from "./Input.style";

interface Props {
  id: string;
  type: "text" | "email" | "password";
}

const Input = ({ type, id }: Props) => {
  return <Styled.Input type={type} id={id} />;
};

export default Input;
