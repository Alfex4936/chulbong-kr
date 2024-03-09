import { ChangeEvent, useState } from "react";

const useInput = (initValue: string) => {
  const [inputValue, setInputValue] = useState(initValue);

  const handleChange = (e: ChangeEvent<HTMLInputElement>) => {
    setInputValue(e.target.value);
  };

  const resetValue = () => {
    setInputValue("");
  };

  return { value: inputValue, onChange: handleChange, reset: resetValue };
};

export default useInput;
