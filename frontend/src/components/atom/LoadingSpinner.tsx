import React from "react";

type Props = {
  size?: "sm" | "md" | "lg";
};

const LoadingSpinner = ({ size = "md" }: Props) => {
  const getSize = () => {
    if (size === "sm") return "h-8 w-8";
    if (size === "md") return "h-12 w-12";
    if (size === "lg") return "h-16 w-16";
  };

  return (
    <div className="flex justify-center items-center">
      <div
        className={`animate-spin rounded-full ${getSize()} border-b border-grey`}
      ></div>
    </div>
  );
};

export default LoadingSpinner;
