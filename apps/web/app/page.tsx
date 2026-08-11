'use client';

import { useState } from 'react';
import { loadEngineWasm, engineVersion } from '@/lib/engine-wasm/wasm';
import { useGameStore } from '@/stores/game-store';

export default function Home() {
  const [status, setStatus] = useState<string>('ยังไม่ได้โหลด engine');
  const [loading, setLoading] = useState(false);
  const setEngineVersion = useGameStore((s) => s.setEngineVersion);

  async function handleLoadWasm() {
    setLoading(true);
    try {
      await loadEngineWasm();
      const v = engineVersion();
      setEngineVersion(v);
      setStatus(`✅ WASM engine พร้อมใช้ (v${v}) — Local/AI mode ทำงานได้`);
    } catch (e) {
      setStatus(`❌ ${(e as Error).message} (ลองรัน \`pnpm build:wasm\` ก่อน)`);
    } finally {
      setLoading(false);
    }
  }

  return (
    <main className="mx-auto flex min-h-screen max-w-3xl flex-col gap-6 px-6 py-12">
      <header>
        <h1 className="text-4xl font-bold text-emerald-400">🎲 Finance Boardgame</h1>
        <p className="mt-2 text-slate-400">
          เกมกระดานการเงินแนว Cashflow — Universal Engine (Go + WASM) + Next.js
        </p>
      </header>

      <section className="rounded-lg border border-slate-700 bg-slate-800/50 p-6">
        <h2 className="mb-3 text-lg font-semibold text-slate-200">
          🔧 ทดสอบ Engine Bridge
        </h2>
        <p className="mb-4 text-sm text-slate-400">
          ปุ่มด้านล่างโหลด engine (Go) เป็น WebAssembly ใน browser เพื่อยืนยันว่า
          bridge ระหว่าง Next.js ↔ WASM ทำงานได้ (foundation ของ Local/AI mode)
        </p>
        <button
          onClick={handleLoadWasm}
          disabled={loading}
          className="rounded-md bg-emerald-600 px-4 py-2 font-medium text-white transition hover:bg-emerald-500 disabled:opacity-50"
        >
          {loading ? 'กำลังโหลด…' : 'โหลด WASM Engine'}
        </button>
        <p className="mt-3 text-sm text-slate-300">{status}</p>
      </section>

      <section className="rounded-lg border border-slate-700 bg-slate-800/50 p-6 text-sm">
        <h2 className="mb-3 text-lg font-semibold text-slate-200">📐 สถาปัตยกรรม</h2>
        <ul className="space-y-1 text-slate-400">
          <li>🎯 <code className="text-emerald-400">packages/engine</code> — Go engine (pure, deterministic)</li>
          <li>🟢 <code className="text-emerald-400">apps/backend</code> — Go service (REST + WS) สำหรับ Online/Async</li>
          <li>⚛️ <code className="text-emerald-400">apps/web</code> — Next.js + PWA (UI และ engine host)</li>
        </ul>
        <p className="mt-4 text-xs text-slate-500">
          รายละเอียดดูใน NOTE.md (Session #2) และ REQUIREMENT.md
        </p>
      </section>
    </main>
  );
}
