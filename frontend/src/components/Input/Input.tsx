import {
  ChangeEvent,
  useRef,
  useState,
  ComponentProps,
  ReactNode,
} from "react";
import Button from "@mui/material/Button";
import * as Styled from "./Input.style";

interface Props extends ComponentProps<"input"> {
  id: string;
  type: "text" | "email" | "password" | "number";
  placeholder: string;
  value: string;
  theme?: "button";
  buttonText?: string | ReactNode;
  onClickFn?: VoidFunction;
  onChange: (e: ChangeEvent<HTMLInputElement>) => void;
}

const Input = ({
  type,
  id,
  placeholder,
  value,
  theme,
  buttonText,
  onClickFn,
  onChange,
  ...props
}: Props) => {
  const [action, setAction] = useState(0);

  const inputRef = useRef<HTMLInputElement | null>(null);

  return (
    <Styled.InputWrap>
      <Styled.Placeholder
        action={action}
        onClick={() => {
          if (inputRef.current) {
            inputRef.current.focus();
          }
        }}
      >
        {placeholder}
      </Styled.Placeholder>
      <Styled.Input
        action={action}
        type={type}
        id={id}
        ref={inputRef}
        value={value}
        onChange={onChange}
        onFocus={() => {
          setAction(1);
        }}
        onBlur={() => {
          if (value === "") setAction(0);
        }}
        {...props}
      />
      {theme === "button" && onClickFn && (
        <Button
          sx={{
            position: "absolute",
            right: "0",
            bottom: "5px",

            height: "26px",

            color: "#fff",
            fontSize: ".6rem",
            backgroundColor: "#333",
            "&:hover": {
              backgroundColor: "#555",
            },
          }}
          onClick={onClickFn}
        >
          {buttonText}
        </Button>
      )}
    </Styled.InputWrap>
  );
};

export default Input;
