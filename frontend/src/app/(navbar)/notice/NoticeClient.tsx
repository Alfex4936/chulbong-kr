"use client";

import useRoadviewStatusStore from "@/store/useRoadviewStatusStore";

const NoticeClient = () => {
  const { open } = useRoadviewStatusStore();

  return <button onClick={open}>fasd</button>;
};

export default NoticeClient;
