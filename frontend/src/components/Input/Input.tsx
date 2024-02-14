import { useRef, useState } from "react";
import * as Styled from "./Input.style";
import useInput from "../../hooks/useInput";

interface Props {
  id: string;
  type: "text" | "email" | "password";
  placeholder: string;
}

const Input = ({ type, id, placeholder }: Props) => {
  const [action, setAction] = useState(0);
  const input = useInput("");

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
        type={type}
        id={id}
        ref={inputRef}
        value={input.value}
        onChange={input.onChange}
        onFocus={() => {
          setAction(1);
        }}
        onBlur={() => {
          if (input.value === "") setAction(0);
        }}
      />
    </Styled.InputWrap>
  );
};

export default Input;
