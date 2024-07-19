import AlertLogin from "@/components/layout/AlertLogin";
import ImageDetail from "@/components/layout/ImageDetail";
import Navigation from "@/components/layout/Navigation";
import MapWrapper from "@/components/map/MapWrapper";
import ChatIdProvider from "@/components/provider/ChatIdProvider";
import PwaAlert from "@/components/provider/PwaAlert";
import RQProvider from "@/components/provider/RQProvider";
import { ThemeProvider } from "@/components/provider/theme-provider";
import { Toaster } from "@/components/ui/toaster";
import type { Metadata, Viewport } from "next";
import localFont from "next/font/local";

import "../globals.css";
import BodyToggleButton from "@/components/layout/body-toggle-button";

declare global {
  interface Window {
    kakao: any;
  }
}

const pretendard = localFont({
  src: [
    {
      path: "../font/Pretendard-Light.woff",
      weight: "300",
      style: "normal",
    },
    {
      path: "../font/Pretendard-Regular.woff",
      weight: "400",
      style: "normal",
    },
    {
      path: "../font/Pretendard-Medium.woff",
      weight: "500",
      style: "normal",
    },
    {
      path: "../font/Pretendard-Bold.woff",
      weight: "700",
      style: "normal",
    },
    {
      path: "../font/Pretendard-ExtraBold.woff",
      weight: "800",
      style: "normal",
    },
  ],
});

export const viewport: Viewport = {
  themeColor: "#222",
};

export const metadata: Metadata = {
  metadataBase: new URL(process.env.NEXT_PUBLIC_URL as string),
  title: "대한민국 철봉 지도",
  keywords: "철봉지도,위치등록,철봉정보,채팅,위치검색,관리,철봉찾기",
  description:
    "대한민국 철봉 지도는 전국 공원의 철봉 위치를 사용자가 직접 등록하고 조회할 수 있는 플랫폼입니다. 가까운 곳에서 철봉 운동을 하고 싶은 분들을 위해, 실시간으로 업데이트되는 철봉 정보를 제공합니다.",
  openGraph: {
    type: "website",
    url: "https://www.k-pullup.com",
    title: "대한민국 철봉 지도",
    description:
      "가까운 곳에서 철봉 위치를 찾고 운동에 참여하세요! 철봉 맵은 전국 공원의 철봉 위치를 사용자가 직접 등록하고 조회할 수 있는 플랫폼입니다.",
    images: "/images/metaimg.webp",
  },
  twitter: {
    card: "summary_large_image",
    title: "대한민국 철봉 지도",
    description:
      "가까운 곳에서 철봉 위치를 찾고 운동에 참여하세요! 철봉 맵은 전국 공원의 철봉 위치를 사용자가 직접 등록하고 조회할 수 있는 플랫폼입니다.",
    images: "/images/metaimg.webp",
  },
  verification: {
    google: "xsTAtA1ny-_9QoSKUsxC7zk_LljW5KBbcWULaNl2gt8",
    other: { naver: "d1ba940a668490789711101918c8b1f7e221a178" },
  },
};

const RootLayout = ({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) => {
  return (
    <html lang="ko">
      <head>
        <link
          rel="manifest"
          href="/manifest.json"
          crossOrigin="use-credentials"
        />
      </head>
      <body
        className={`${pretendard.className} overflow-hidden flex text-black dark:text-grey mo:flex-col-reverse h-dvh mo:bg-neutral-800`}
      >
        <ThemeProvider
          attribute="class"
          defaultTheme="dark"
          enableSystem
          disableTransitionOnChange
        >
          <RQProvider>
            <ChatIdProvider>
              <Navigation />
              {children}
              <BodyToggleButton />
              <ImageDetail />
              <PwaAlert />
              <Toaster />
              <MapWrapper />
              <AlertLogin />
            </ChatIdProvider>
          </RQProvider>
        </ThemeProvider>
      </body>
    </html>
  );
};

export default RootLayout;
