import RQProvider from "@/components/provider/RQProvider";
import type { Metadata } from "next";
import { Nanum_Gothic } from "next/font/google";
import Head from "next/head";
import "../globals.css";

const nanum = Nanum_Gothic({
  subsets: ["latin"],
  weight: ["400", "700", "800"],
  display: "swap",
});
export const metadata: Metadata = {
  title: "대한민국 철봉 지도",
  description:
    "대한민국 철봉 지도는 전국 공원의 철봉 위치를 사용자가 직접 등록하고 조회할 수 있는 플랫폼입니다. 가까운 곳에서 철봉 운동을 하고 싶은 분들을 위해, 실시간으로 업데이트되는 철봉 정보를 제공합니다.",
};

const RootLayout = ({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) => {
  return (
    <html lang="ko">
      <Head>
        <meta name="viewport" content="initial-scale=1, width=device-width" />
        <meta
          name="keywords"
          content="철봉지도,위치등록,철봉정보,채팅,위치검색,관리,철봉찾기"
        />
        <meta property="og:type" content="website" />
        <meta property="og:url" content="https://www.k-pullup.com" />
        <meta property="og:title" content="대한민국 철봉 지도" />
        <meta
          property="og:description"
          content="가까운 곳에서 철봉 위치를 찾고 운동에 참여하세요! 철봉 맵은 전국 공원의 철봉 위치를 사용자가 직접 등록하고 조회할 수 있는 플랫폼입니다."
        />
        <meta property="og:image" content="/images/metaimg.webp" />
        <meta property="twitter:card" content="summary_large_image" />
        <meta property="twitter:url" content="https://www.k-pullup.com" />
        <meta property="twitter:title" content="대한민국 철봉 지도" />
        <meta
          property="twitter:description"
          content="철봉 맵은 전국 공원의 철봉 위치를 사용자가 직접 등록하고 조회할 수 있는 플랫폼입니다."
        />
        <meta property="twitter:image" content="/images/metaimg.webp" />
        <meta
          name="google-site-verification"
          content="xsTAtA1ny-_9QoSKUsxC7zk_LljW5KBbcWULaNl2gt8"
        />
        <meta
          name="naver-site-verification"
          content="d1ba940a668490789711101918c8b1f7e221a178"
        />
      </Head>
      <body
        className={`${nanum.className} bg-black-gradient-2 h-screen w-screen text-grey`}
      >
        <RQProvider>{children}</RQProvider>
      </body>
    </html>
  );
};

export default RootLayout;
