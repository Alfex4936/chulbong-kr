import { cn } from "@/lib/utils";
import { ComponentProps } from "react";

interface Props extends ComponentProps<"button"> {
  icon: React.ReactNode;
  className?: string;
  numberState?: number;
  theme?: "dark" | "light";
}

const IconButton = ({
  icon,
  className,
  numberState,
  theme = "light",
  ...props
}: Props) => {
  return (
    <button
      {...props}
      className={cn(
        `w-8 h-8 absolute p-1 z-20 rounded-sm ${
          theme === "dark"
            ? "bg-black-light-2 hover:bg-grey-dark-1"
            : "bg-white-tp-light hover:bg-white-tp-dark"
        }`,
        className
      )}
    >
      {typeof numberState === "number" && (
        <div className="absolute -top-1 -right-1 text-[10px] px-2 py-[1px] bg-red rounded-lg">
          {numberState}
        </div>
      )}

      {icon}
    </button>
  );
};

export default IconButton;
