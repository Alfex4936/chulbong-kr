/** @type {import('next').NextConfig} */

import withPWA from "next-pwa";

const nextPWA = withPWA({
  dest: "public",
});

const nextConfig = {
  reactStrictMode: false,
  async rewrites() {
    return [
      {
        source: "/api/v1/:path*",
        destination: "https://api.k-pullup.com/api/v1/:path*",
      },
    ];
  },
  images: {
    domains: ["chulbong-kr.s3.amazonaws.com"],
  },
  crossOrigin: "use-credentials",
};

export default nextConfig;

// import { MetadataRoute } from "next";

// const manifest = (): MetadataRoute.Manifest => {
//   return {
//     theme_color: "#222",
//     background_color: "#222",
//     display: "standalone",
//     scope: "/",
//     start_url: "/",
//     name: "철봉 지도",
//     short_name: "철봉 지도",
//     description: "주변 철봉 위치를 확인하세요",
//     lang: "ko-KR",
//     icons: [
//       {
//         src: "/logo192.png",
//         sizes: "192x192",
//         type: "image/png",
//         purpose: "maskable",
//       },
//       {
//         src: "/logo512.png",
//         sizes: "512x512",
//         type: "image/png",
//       },
//     ],
//   };
// };
// export default manifest;
