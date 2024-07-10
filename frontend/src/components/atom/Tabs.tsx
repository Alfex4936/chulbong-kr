"use client";

import { useEffect, useState } from "react";
import useTabStore from "@/store/useTabStore";

type Props = {
  title: string;
  tabs: string[];
  children: React.ReactNode;
};

const Tabs = ({ title, tabs, children }: Props) => {
  const { setCurTab, disableIndex } = useTabStore();
  const [tabIndex, setTabIndex] = useState(0);

  useEffect(() => {
    setCurTab(tabs[0]);
  }, [tabs[0]]);

  return (
    <div className="relative">
      <div className="flex mb-3">
        <div className="text-xl mr-1">{title}</div>
        <div className="flex items-end">
          {tabs.map((tab, index) => {
            return (
              <button
                key={tab}
                className={`text-xs ${
                  tabIndex === index ? "text-grey" : "text-grey-dark"
                } px-1 ${tabIndex === index && "underline"}`}
                onClick={() => {
                  setTabIndex(index);
                  setCurTab(tab);
                }}
                disabled={disableIndex === index}
              >
                {tab}
              </button>
            );
          })}
        </div>
      </div>
      {children}
    </div>
  );
};

export default Tabs;
