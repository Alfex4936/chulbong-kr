import usePageLoadingStore from "@/store/usePageLoadingStore";
import { useEffect, useState } from "react";

const PageLoadingBar = () => {
  const { isLoading } = usePageLoadingStore();
  const [width, setWidth] = useState("0%");
  const [visible, setVisible] = useState(false);

  useEffect(() => {
    if (isLoading) {
      setVisible(true);
      setWidth("30%");
    } else {
      setWidth("100%");
      const time = setTimeout(() => {
        setVisible(false);
      }, 300);

      return () => {
        clearTimeout(time);
      };
    }
  }, [isLoading]);

  return (
    visible && (
      <div
        className="fixed top-0 left-0 h-[2px] bg-grey-dark-1 z-50 transition-all duration-200"
        style={{ width }}
      ></div>
    )
  );
};

export default PageLoadingBar;
