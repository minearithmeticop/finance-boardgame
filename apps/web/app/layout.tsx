import type { Metadata, Viewport } from 'next';
import './globals.css';

export const metadata: Metadata = {
  title: 'Finance Boardgame',
  description:
    'เกมกระดานการเงินแนว Cashflow — เรียนรู้ Financial Literacy ผ่านการบริหาร Cashflow',
  applicationName: 'Finance Boardgame',
  appleWebApp: { capable: true, statusBarStyle: 'default', title: 'Finance Boardgame' },
};

export const viewport: Viewport = {
  themeColor: '#0f172a',
  width: 'device-width',
  initialScale: 1,
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="th">
      <body>{children}</body>
    </html>
  );
}
