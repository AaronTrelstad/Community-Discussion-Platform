import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  async rewrites() {
    return [
        {
            source: "/api/auth/:path*",
            destination: "http://localhost:8080/auth/:path*",
        },
        {
            source: "/api/communities/:path*",
            destination: "http://localhost:8081/communities/:path*",
        },
    ];
},
};

export default nextConfig;
