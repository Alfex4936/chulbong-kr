import { MetadataRoute } from "next";

const manifest = (): MetadataRoute.Manifest => {
  return {
    theme_color: "#222",
    background_color: "#222",
    display: "standalone",
    scope: "/",
    start_url: "/",
    name: "철봉 지도",
    short_name: "철봉 지도",
    description: "주변 철봉 위치를 확인하세요",
    lang: "ko-KR",
    icons: [
      {
        src: "/2.png",
        sizes: "192x192",
        type: "image/png",
      },
    ],
  };
};
export default manifest;
