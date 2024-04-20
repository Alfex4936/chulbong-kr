import { ComponentProps } from "react";

interface Props extends ComponentProps<"button"> {
  icon: React.ReactNode;
  selected: boolean;
  text?: string;
}

const IconButton = ({ icon, selected, text, ...props }: Props) => {
  return (
    <button
      className="flex justify-center items-center web:mb-3 mo:w-full"
      {...props}
    >
      <div className="flex flex-col justify-center items-center w-12 h-16 cursor-pointer rounded-sm mo:w-[44px] mo:h-[50px]">
        <div className="ml-[1px]">{icon}</div>
        {text && (
          <div
            className={` text-sm mp:font-xs mo:text-[10px] web:mt-2 text-grey`}
          >
            {text}
          </div>
        )}
      </div>
    </button>
  );
};

export default IconButton;
