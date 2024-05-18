import { cn } from "@/utils/cn";
import formatDigitNumber from "@/utils/formatDigitNumber";
import { useEffect, useState } from "react";

interface Props {
  start: boolean;
  setStart: React.Dispatch<React.SetStateAction<boolean>>;
  initTime?: number;
  className?: string;
}
const Count = ({ start, setStart, initTime = 300, className }: Props) => {
  const [time, setTime] = useState(initTime);
  const [fontColor, setFontColor] = useState("#6767ff");

  useEffect(() => {
    let interval: NodeJS.Timeout;
    if (start) {
      setFontColor("#6767ff");
      interval = setInterval(() => {
        setTime((prev) => prev - 1);
      }, 1000);
    }

    if (time === 0) {
      setFontColor("red");
      setTime(initTime);
      setStart(false);
    }

    if (start === false) {
      setTime(initTime);
    }

    return () => {
      clearInterval(interval);
    };
  }, [time, start]);

  const formatTime = (timer: number): string => {
    const minutes = Math.floor(timer / 60);
    const seconds = timer % 60;
    return `${formatDigitNumber(minutes)} : ${formatDigitNumber(seconds)}`;
  };

  return (
    <div style={{ color: fontColor }} className={cn(className)}>
      {formatTime(time)}
    </div>
  );
};

export default Count;
