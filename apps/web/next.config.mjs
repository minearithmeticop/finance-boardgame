/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  // ให้ Next ให้บริการ /public/wasm/*.wasm เป็น static asset ได้
  // (ไม่ต้องตั้งค่าพิเศษ — ไฟล์ใน public/ ถูก serve ที่ root path อยู่แล้ว)
};

export default nextConfig;
