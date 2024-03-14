import Tooltip from "@mui/material/Tooltip";
import IconButton from "@mui/material/IconButton";
import { useEffect, useState } from "react";

export interface FloatingProps {
  text: string | React.ReactNode;
  top?: number | "center" | string;
  right?: number | "center" | string;
  bottom?: number | "center" | string;
  left?: number | "center" | string;
  width?: number | string;
  height?: number | string;
  zIndex?: number;
  shape?: "circle" | "square";
  size?: "small" | "large" | "medium";
  tooltip?: string;
  onClickFn?: VoidFunction;
}

const FloatingButton = ({
  text,
  top,
  right,
  bottom,
  left,
  width,
  height,
  zIndex = 10,
  shape = "square",
  size = "medium",
  tooltip,
  onClickFn,
}: FloatingProps) => {
  const [textSize, setTextSize] = useState("");
  const [boxSize, setBoxSize] = useState("");

  useEffect(() => {
    if (size === "small") {
      setTextSize("12px");
      setBoxSize("50px");
    }
    if (size === "medium") {
      setTextSize("16px");
      setBoxSize("40px");
    }
    if (size === "large") {
      setTextSize("20px");
      setBoxSize("30px");
    }
  }, [size]);

  return (
    <Tooltip title={tooltip} arrow disableInteractive>
      <IconButton
        // size={size}
        sx={{
          position: "absolute",
          top: top === "center" ? "50%" : `${top}px`,
          right: right === "center" ? "50%" : `${right}px`,
          bottom: bottom === "center" ? "50%" : `${bottom}px`,
          left: left === "center" ? "50%" : `${left}px`,

          transform:
            top === "center" || bottom === "center"
              ? "translateY(-50%)"
              : right === "center" || left === "center"
              ? "translateX(-50%)"
              : 0,

          width: width ? `${width}px` : boxSize,
          height: height ? `${height}px` : boxSize,

          fontSize: textSize,

          boxShadow:
            "rgba(50, 50, 93, 0.25) 0px 2px 5px -1px, rgba(0, 0, 0, 0.3) 0px 1px 3px -1px",

          backgroundColor: "#fff",
          color: "#000",

          zIndex: zIndex,

          borderRadius: shape === "square" ? "3px" : "50%",

          "&:hover": {
            backgroundColor: "#888",
            color: "#fff",
          },
        }}
        onClick={onClickFn}
      >
        <div>{text}</div>
      </IconButton>
    </Tooltip>
  );
};

export default FloatingButton;
