"use client";

import useTabStore from "@/store/useTabStore";

const RankingResult = () => {
  const { curTab } = useTabStore();

  return <div>{curTab}</div>;
};

export default RankingResult;
