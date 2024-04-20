import React from "react";

type Props = {
  size?: "xs" | "sm" | "md" | "lg";
  color?: "white" | "black";
};

const LoadingSpinner = ({ size = "md", color = "white" }: Props) => {
  const getSize = () => {
    if (size === "xs") return "h-4 w-4";
    if (size === "sm") return "h-8 w-8";
    if (size === "md") return "h-12 w-12";
    if (size === "lg") return "h-16 w-16";
  };

  const getColor = () => {
    if (color === "white") return "border-grey";
    if (color === "black") return "border-black";
  };

  return (
    <div className="flex justify-center items-center">
      <div
        className={`animate-spin rounded-full ${getSize()} border-b ${getColor()}`}
      ></div>
    </div>
  );
};

export default LoadingSpinner;
