"use client";

import Image from "next/image";
import { useEffect, useState } from "react";
import { LuUpload } from "react-icons/lu";

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

type DeviceType = "IOS Chrome" | "IOS Safari" | "IOS" | "Android" | "Web";

interface NavigatorStandalone extends Navigator {
  standalone?: boolean;
}

const isRunningStandalone = (): boolean => {
  const navigatorWithStandalone = navigator as NavigatorStandalone;

  return (
    window.matchMedia("(display-mode: standalone)").matches ||
    navigatorWithStandalone.standalone === true
  );
};

const PwaAlert = () => {
  const [alert, setAlert] = useState(false);
  const [prompt, setPrompt] = useState<BeforeInstallPromptEvent | null>(null);

  const [isApp, setIsApp] = useState(false);
  const [isMobile, setIsMobile] = useState(false);

  const [downInfo, setDownInfo] = useState(false);

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

    if (isRunningStandalone()) {
      setIsApp(true);
    } else {
      setIsApp(false);
    }

    if (getDeviceType() === ("IOS" || "Android")) {
      setIsMobile(true);
    } else {
      setIsMobile(false);
    }

    handleResize();

    window.addEventListener("resize", handleResize);
    window.addEventListener("beforeinstallprompt", handlePrompt);

    return () => {
      window.removeEventListener("resize", handleResize);
      window.removeEventListener("beforeinstallprompt", handlePrompt);
    };
  }, []);

  const getDeviceType = (): DeviceType => {
    const userAgent = navigator.userAgent;

    if (/iPad|iPhone|iPod/.test(userAgent)) {
      if (/CriOS/.test(userAgent)) {
        return "IOS Chrome";
      }
      if (/Safari/.test(userAgent) && !/CriOS/.test(userAgent)) {
        return "IOS Safari";
      }
      return "IOS";
    }

    if (/android/i.test(userAgent)) {
      return "Android";
    }

    return "Web";
  };

  const handleInstallClick = () => {
    if (prompt) {
      prompt.prompt();
      prompt.userChoice.then(() => {
        setPrompt(null);
      });
    } else {
      if (isMobile) {
        setDownInfo(true);
        return;
      }
      setAlert(false);
    }
  };

  if (!alert || isApp) return null;

  return (
    <div className="absolute top-0 left-0 w-dvw h-dvh bg-white-tp-light z-[900]">
      {downInfo ? (
        <>
          {getDeviceType() === "IOS" || "IOS Chrome" ? (
            <>
              <div className="absolute left-1/2 -translate-x-1/2 top-28 w-[90%] bg-black-light-2 z-[1000] p-4 rounded-md">
                <button
                  className="absolute top-1 right-2"
                  onClick={() => setAlert(false)}
                >
                  X
                </button>
                <div className="mb-3 text-center">
                  화면 상단에 다운로드 아이콘을 클릭하여
                  <br /> 홈 화면에 추가해 주세요!
                </div>
                <div className="text-4xl rounded-full w-14 h-14 flex items-center justify-center mx-auto bg-grey-dark-1">
                  <LuUpload />
                </div>
              </div>
              <div className="absolute top-1 right-1">
                <Image
                  src={"/arrowcu.png"}
                  width={40}
                  height={100}
                  alt="arrow"
                />
              </div>
            </>
          ) : (
            <>
              <div className="absolute left-1/2 -translate-x-1/2 top-28 w-[90%] bg-black-light-2 z-[1000] p-4 rounded-md">
                <button
                  className="absolute top-1 right-2"
                  onClick={() => setAlert(false)}
                >
                  X
                </button>
                <div className="mb-3 text-center">
                  화면 상단에 다운로드 아이콘을 클릭하여
                  <br /> 홈 화면에 추가해 주세요!
                </div>
                <div className="text-4xl rounded-full w-14 h-14 flex items-center justify-center mx-auto bg-grey-dark-1">
                  <LuUpload />
                </div>
              </div>
              <div className="absolute bottom-1 left-1/2 -translate-x-1/2">
                <Image
                  src={"/arrowcd.png"}
                  width={40}
                  height={100}
                  alt="arrow"
                />
              </div>
            </>
          )}
        </>
      ) : (
        <>
          <div className="absolute left-1/2 -translate-x-1/2 bottom-10 w-[90%] bg-black-light-2 z-[1000] p-4 rounded-md">
            <div className="text-base mb-2">
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
        </>
      )}
    </div>
  );
};

export default PwaAlert;
