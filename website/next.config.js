/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  swcMinify: true,
  images: {
    domains: ['github.com'],
  },
  async redirects() {
    return [
      {
        source: '/download',
        destination: 'https://github.com/Nom-nom-hub/gotunnel/releases',
        permanent: true,
      },
    ]
  },
}

module.exports = nextConfig 