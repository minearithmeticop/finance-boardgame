import type { Metadata, Viewport } from 'next';
import Link from 'next/link';
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
      <body>
        <nav className="border-b border-slate-800 bg-slate-900/80 backdrop-blur">
          <div className="mx-auto flex max-w-5xl items-center gap-6 px-6 py-3 text-sm">
            <Link href="/" className="font-semibold text-emerald-400">
              🎲 Finance Boardgame
            </Link>
            <Link
              href="/glossary"
              className="text-slate-300 transition hover:text-emerald-400"
            >
              📚 คัมภีร์การเงิน
            </Link>
          </div>
        </nav>
        {children}
      </body>
    </html>
  );
}
