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

export default nextPWA(nextConfig);
