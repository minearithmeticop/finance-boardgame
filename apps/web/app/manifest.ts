import type { MetadataRoute } from 'next';

// PWA manifest — Next.js จะ generate เป็น /manifest.webmanifest อัตโนมัติ
export default function manifest(): MetadataRoute.Manifest {
  return {
    name: 'Finance Boardgame',
    short_name: 'FinanceGame',
    description: 'เกมกระดานการเงินแนว Cashflow',
    start_url: '/',
    display: 'standalone',
    background_color: '#0f172a',
    theme_color: '#0f172a',
    icons: [
      { src: '/icon.svg', sizes: 'any', type: 'image/svg+xml' },
    ],
  };
}
