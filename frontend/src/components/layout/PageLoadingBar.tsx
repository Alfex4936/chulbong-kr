import usePageLoadingStore from "@/store/usePageLoadingStore";
import { useEffect, useState } from "react";

const PageLoadingBar = () => {
  const { isLoading, visible, setVisible } = usePageLoadingStore();
  const [width, setWidth] = useState("30%");

  useEffect(() => {
    if (isLoading) {
      setWidth("30%");
      setVisible(true);
    } else {
      setWidth("100%");
    }
  }, [isLoading]);

  return (
    visible && (
      <div
        className="fixed top-0 left-0 h-[3px] w-1/3 bg-grey-dark-1 transition-all duration-200 z-[500]"
        style={{ width }}
      ></div>
    )
  );
};

export default PageLoadingBar;
