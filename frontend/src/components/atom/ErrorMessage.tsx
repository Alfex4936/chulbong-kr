type Props = {
  text: string;
};

const ErrorMessage = ({ text }: Props) => {
  return <div className="text-xs text-red">{text}</div>;
};

export default ErrorMessage;
