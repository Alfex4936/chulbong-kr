"use client";

import useChatIdStore from "@/store/useChatIdStore";
import { useEffect } from "react";

const ChatIdProvider = ({ children }: { children: React.ReactNode }) => {
  const cidState = useChatIdStore();

  useEffect(() => {
    const setId = (e?: StorageEvent) => {
      if (e) {
        if (e.key === "cid") {
          cidState.setId();
        }
      } else {
        if (cidState.cid) return;
        cidState.setId();
      }
    };

    setId();

    window.addEventListener("storage", setId);

    return () => {
      window.removeEventListener("storage", setId);
    };
  }, []);
  return <>{children}</>;
};

export default ChatIdProvider;
