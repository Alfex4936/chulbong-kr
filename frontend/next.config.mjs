/** @type {import('next').NextConfig} */

// import withPWA from "next-pwa";

// const nextPWA = withPWA({
//   dest: "public",
// });

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
  async headers() {
    return [
      {
        source: "/manifest.webmanifest",
        headers: [
          {
            key: "Access-Control-Allow-Origin",
            value: "*",
          },
          {
            key: "Access-Control-Allow-Methods",
            value: "GET,OPTIONS",
          },
          {
            key: "Access-Control-Allow-Headers",
            value: "Content-Type",
          },
        ],
      },
    ];
  },
  images: {
    domains: ["chulbong-kr.s3.amazonaws.com"],
  },
};

export default nextConfig;
