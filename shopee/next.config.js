/{id}{id} @type {import('next').NextConfig} {id}/
const nextConfig = {
    reactStrictMode: true,
    images: {
        domains: ['res.cloudinary.com'],
    },
};

module.exports = nextConfig;
