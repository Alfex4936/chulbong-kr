/** @type {import('next').NextConfig} */

const nextConfig = {
  crossOrigin: "use-credentials",
  reactStrictMode: false,
  async rewrites() {
    return [
      {
        source: "/api/v1/:path*",
        destination: "https://api.k-pullup.com/api/v1/:path*",
      },
    ];
  },
};

export default nextConfig;
