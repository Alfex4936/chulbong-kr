import { cn } from "@/utils/cn";

type Props = { title: string; subTitle?: string; className?: string };

const Heading = ({ title, subTitle, className }: Props) => {
  return (
    <div
      className={cn(
        "flex flex-col items-center justify-center font-medium text-2xl text-center h-24 mo:text-lg",
        className
      )}
    >
      <div>{title}</div>
      {subTitle && (
        <div className="text-sm text-grey-dark mo:text-xs">{subTitle}</div>
      )}
    </div>
  );
};

export default Heading;
