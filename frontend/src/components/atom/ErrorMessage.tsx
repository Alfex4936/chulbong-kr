import { cn } from "@/lib/utils";

interface Props {
  text?: string;
  className?: string;
}
const ErrorMessage = ({ text, className }: Props) => {
  return <div className={cn("text-xs text-red", className)}>{text}</div>;
};

export default ErrorMessage;
