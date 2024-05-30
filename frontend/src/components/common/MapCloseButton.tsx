import { useEffect, useState } from "react";

const MapCloseButton = () => {
  const [isAnimated, setIsAnimated] = useState(false);

  useEffect(() => {
    // 컴포넌트가 마운트된 후 애니메이션을 시작합니다.
    const timer = setTimeout(() => {
      setIsAnimated(true);
    }, 100);
    return () => clearTimeout(timer); // 컴포넌트가 언마운트될 때 타이머를 정리합니다.
  }, []);

  return (
    <div className="relative w-16 h-16">
      <div
        className={`absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-1/3 h-[1.5px] bg-white transform origin-center transition-transform duration-500 ease-in-out ${
          isAnimated ? "rotate-45" : ""
        }`}
      ></div>
      <div
        className={`absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-1/3 h-[1.5px] bg-white transform origin-center transition-transform duration-500 ease-in-out ${
          isAnimated ? "-rotate-45" : ""
        }`}
      ></div>
    </div>
  );
};

export default MapCloseButton;
