import { ComponentProps } from "react";

interface Props extends ComponentProps<"button"> {
  icon: React.ReactNode;
  top?: number;
  left?: number;
  right?: number;
  bottom?: number;
  numberState?: number;
}

const IconButton = ({
  icon,
  top,
  left,
  right,
  bottom,
  numberState,
  ...props
}: Props) => {
  return (
    <button
      {...props}
      className="w-8 h-8 absolute p-1 z-20 rounded-sm bg-white-tp-light hover:bg-white-tp-dark"
      style={{ top, left, right, bottom }}
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
