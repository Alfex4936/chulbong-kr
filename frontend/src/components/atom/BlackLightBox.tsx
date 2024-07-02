import { cn } from "@/lib/utils";

type Props = {
  children: React.ReactNode;
  center?: boolean;
  className?: string;
};

const BlackLightBox = ({ children, center = false, className }: Props) => {
  return (
    <div
      className={cn(
        `bg-black-light-2 mx-auto w-[90%] p-4 rounded-md ${
          center ? "text-center" : "text-left"
        }`,
        className
      )}
    >
      {children}
    </div>
  );
};

export default BlackLightBox;
