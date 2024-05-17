import { ChangeEvent, useState } from "react";

const useInput = (initValue: string) => {
  const [value, setValue] = useState(initValue);

  const handleChange = (
    e: ChangeEvent<HTMLInputElement> | ChangeEvent<HTMLTextAreaElement>
  ) => {
    setValue(e.target.value);
  };

  const resetValue = () => {
    setValue("");
  };

  return { value, handleChange, resetValue };
};

export default useInput;
