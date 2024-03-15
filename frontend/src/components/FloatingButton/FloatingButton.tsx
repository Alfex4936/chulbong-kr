import IconButton from "@mui/material/IconButton";
import Tooltip from "@mui/material/Tooltip";

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

          fontSize:
            size === "small" ? "12px" : size === "medium" ? "16px" : "20px",
          width: width
            ? width
            : size === "small"
            ? "50px"
            : size === "medium"
            ? "40px"
            : "30px",
          height: height
            ? height
            : size === "small"
            ? "50px"
            : size === "medium"
            ? "40px"
            : "30px",

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
