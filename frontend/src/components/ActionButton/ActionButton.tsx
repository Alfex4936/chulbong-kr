import { ButtonProps } from "@mui/material";
import Button from "@mui/material/Button";

interface Props extends ButtonProps {
  bg?: "black" | "gray";
}

const ActionButton = ({ bg, ...props }: Props) => {
  if (!bg) {
    return <Button>{props.children}</Button>;
  }

  return (
    <Button
      sx={{
        color: bg === "black" ? "#fff" : "#000",
        width: "100%",
        backgroundColor: bg === "black" ? "#333" : "#ccc",
        margin: "1rem 0",
        "&:hover": {
          backgroundColor: bg === "black" ? "#555" : "#eee",
        },
      }}
      {...props}
    >
      {props.children}
    </Button>
  );
};

export default ActionButton;
