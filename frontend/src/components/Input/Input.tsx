interface Props {
  id: string;
  type: "text" | "email" | "password";
}

const Input = ({ type, id }: Props) => {
  return <input type={type} id={id} />;
};

export default Input;
