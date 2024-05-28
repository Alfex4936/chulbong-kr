"use client";

import { useEffect, useState } from "react";

declare global {
  interface WindowEventMap {
    beforeinstallprompt: BeforeInstallPromptEvent;
  }

  interface BeforeInstallPromptEvent extends Event {
    readonly platforms: string[];
    readonly userChoice: Promise<{ outcome: "accepted" | "dismissed" }>;
    prompt(): Promise<void>;
  }
}

const PwaAlert = () => {
  const [alert, setAlert] = useState(false);
  const [prompt, setPrompt] = useState<BeforeInstallPromptEvent | null>(null);

  useEffect(() => {
    const handlePrompt = (e: BeforeInstallPromptEvent) => {
      e.preventDefault();
      setPrompt(e);
    };

    const handleResize = () => {
      if (window.innerWidth <= 540) {
        setAlert(true);
      } else {
        setAlert(false);
      }
    };

    handleResize();

    window.addEventListener("resize", handleResize);
    window.addEventListener("beforeinstallprompt", handlePrompt);

    return () => {
      window.removeEventListener("resize", handleResize);
      window.removeEventListener("beforeinstallprompt", handlePrompt);
    };
  }, []);

  const handleInstallClick = () => {
    if (prompt) {
      prompt.prompt();
      prompt.userChoice.then(() => {
        setPrompt(null);
      });
    } else {
      setAlert(false);
    }
  };

  if (!alert) return null;

  return (
    <div className="absolute top-0 left-0 w-dvw h-dvh bg-white-tp-light z-[900]">
      <div className="absolute left-1/2 -translate-x-1/2 bottom-10 w-[90%] bg-black-light-2 z-[1000] p-4 rounded-md">
        <div className="text-lg mb-2">
          홈 화면에 철봉 앱을 추가하고 <br /> 편하게 사용하세요.
        </div>
        <button
          className="bg-grey-dark-1 w-full p-2 rounded-md mb-2"
          onClick={handleInstallClick}
        >
          설치하고 앱으로 보기
        </button>
        <button
          className="text-sm underline text-grey-dark text-center w-full"
          onClick={() => setAlert(false)}
        >
          웹으로 계속 보기
        </button>
      </div>
    </div>
  );
};

export default PwaAlert;
