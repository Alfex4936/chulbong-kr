import { ChangeEvent, useRef, useState } from "react";
import * as Styled from "./Input.style";

interface Props {
  id: string;
  type: "text" | "email" | "password";
  placeholder: string;
  value: string;
  onChange: (e: ChangeEvent<HTMLInputElement>) => void;
}

const Input = ({ type, id, placeholder, value, onChange }: Props) => {
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
      />
    </Styled.InputWrap>
  );
};

export default Input;
