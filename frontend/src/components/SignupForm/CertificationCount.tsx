import { useEffect, useState } from "react";
import formatDigitNumber from "../../utils/formatDigitNumber";

interface Props {
  start: boolean;
  setStart: React.Dispatch<React.SetStateAction<boolean>>;
}

const CertificationCount = ({ start, setStart }: Props) => {
  const [time, setTime] = useState(300);
  const [fontColor, setFontColor] = useState("#6767ff");

  useEffect(() => {
    let interval: number | undefined;
    if (start) {
      setFontColor("#6767ff");
      interval = setInterval(() => {
        setTime((prev) => prev - 1);
      }, 1000);
    }

    if (time === 0) {
      setFontColor("red");
      setTime(300);
      setStart(false);
    }

    if (start === false) {
      setTime(300);
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

  return <div style={{ color: fontColor }}>{formatTime(time)}</div>;
};

export default CertificationCount;
