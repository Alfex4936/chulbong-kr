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
      <div
        className="flex flex-col justify-center items-center w-16 h-16 cursor-pointer bg-grey rounded-sm mo:w-[44px] mo:h-[50px]"
        style={{
          backgroundColor: selected ? "#f0f0f0" : "#222222",
        }}
      >
        <div className="ml-[1px]">{icon}</div>
        {text && (
          <div
            className={`${
              selected ? "text-black-light" : "text-grey"
            } text-sm mp:font-xs mo:text-[10px] web:mt-2`}
          >
            {text}
          </div>
        )}
      </div>
    </button>
  );
};

export default IconButton;
