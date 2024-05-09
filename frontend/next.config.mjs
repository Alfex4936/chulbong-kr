/** @type {import('next').NextConfig} */
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
};

export default nextConfig;
