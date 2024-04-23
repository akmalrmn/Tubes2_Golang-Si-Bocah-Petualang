/** @type {import('next').NextConfig} */
const nextConfig = {
  env: {
    WIKIPEDIA_API_URL: process.env.WIKIPEDIA_API_URL,
    BACKEND_API_URL: process.env.BACKEND_API_URL,
  },
};

export default nextConfig;
