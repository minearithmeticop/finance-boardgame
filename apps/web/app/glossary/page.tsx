'use client';

import { useMemo, useState } from 'react';
import { glossary, type GlossaryCategory } from '@/data/glossary';

type Filter = GlossaryCategory | 'all';

const CATEGORIES: { id: Filter; label: string }[] = [
  { id: 'all', label: 'ทั้งหมด' },
  { id: 'income', label: 'รายได้' },
  { id: 'statement', label: 'งบการเงิน' },
  { id: 'expense-tax', label: 'รายจ่ายและภาษี' },
  { id: 'asset-liability', label: 'สินทรัพย์และหนี้สิน' },
  { id: 'game', label: 'กลไกเกม' },
];

const CATEGORY_STYLE: Record<GlossaryCategory, string> = {
  income: 'bg-emerald-500/15 text-emerald-300 border-emerald-500/30',
  statement: 'bg-sky-500/15 text-sky-300 border-sky-500/30',
  'expense-tax': 'bg-amber-500/15 text-amber-300 border-amber-500/30',
  'asset-liability': 'bg-violet-500/15 text-violet-300 border-violet-500/30',
  game: 'bg-rose-500/15 text-rose-300 border-rose-500/30',
};

const CATEGORY_LABEL: Record<GlossaryCategory, string> = {
  income: 'รายได้',
  statement: 'งบการเงิน',
  'expense-tax': 'รายจ่ายและภาษี',
  'asset-liability': 'สินทรัพย์และหนี้สิน',
  game: 'กลไกเกม',
};

export default function GlossaryPage() {
  const [query, setQuery] = useState('');
  const [cat, setCat] = useState<Filter>('all');

  const filtered = useMemo(() => {
    const q = query.trim().toLowerCase();
    return glossary.filter((t) => {
      if (cat !== 'all' && t.category !== cat) return false;
      if (!q) return true;
      return (
        t.term.toLowerCase().includes(q) ||
        t.th.toLowerCase().includes(q) ||
        t.short.toLowerCase().includes(q) ||
        t.definition.toLowerCase().includes(q)
      );
    });
  }, [query, cat]);

  return (
    <main className="mx-auto max-w-5xl px-6 py-10">
      <header className="mb-8">
        <h1 className="text-3xl font-bold text-emerald-400">📚 คัมภีร์การเงิน</h1>
        <p className="mt-2 max-w-2xl text-slate-400">
          รวมคำศัพท์และแนวคิดทางการเงินในเกม — อธิบายภาษาง่าย พร้อมที่มาที่ไปและตัวอย่าง
          ทั้งผู้เล่นและทีมสร้างเกมเปิดดูเข้าใจตรงกัน
        </p>
      </header>

      {/* ค้นหา + กรอง */}
      <div className="mb-6 flex flex-col gap-3">
        <input
          type="search"
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          placeholder="ค้นหาคำศัพท์ เช่น ภาษี, passive income, ประกันสังคม..."
          className="w-full rounded-md border border-slate-700 bg-slate-800/60 px-4 py-2 text-slate-100 outline-none placeholder:text-slate-500 focus:border-emerald-500"
        />
        <div className="flex flex-wrap gap-2">
          {CATEGORIES.map((c) => (
            <button
              key={c.id}
              onClick={() => setCat(c.id)}
              className={`rounded-full border px-3 py-1 text-xs transition ${
                cat === c.id
                  ? 'border-emerald-500 bg-emerald-500/20 text-emerald-300'
                  : 'border-slate-700 text-slate-400 hover:border-slate-500'
              }`}
            >
              {c.label}
            </button>
          ))}
        </div>
      </div>

      <p className="mb-4 text-sm text-slate-500">พบ {filtered.length} คำ</p>

      {/* รายการคำศัพท์ */}
      <div className="grid gap-4 sm:grid-cols-2">
        {filtered.map((t) => (
          <article
            key={t.id}
            className="flex flex-col rounded-lg border border-slate-700 bg-slate-800/40 p-5"
          >
            <div className="mb-2 flex items-start justify-between gap-3">
              <div>
                <h2 className="text-lg font-semibold text-slate-100">{t.term}</h2>
                <p className="text-sm text-slate-400">{t.th}</p>
              </div>
              <span
                className={`shrink-0 rounded-full border px-2 py-0.5 text-xs ${CATEGORY_STYLE[t.category]}`}
              >
                {CATEGORY_LABEL[t.category]}
              </span>
            </div>

            <p className="mb-1 text-slate-300">{t.short}</p>
            <p className="text-sm text-slate-400">{t.definition}</p>

            <details className="mt-3 text-sm">
              <summary className="cursor-pointer text-emerald-400 hover:text-emerald-300">
                ดูรายละเอียดเพิ่มเติม
              </summary>
              <div className="mt-2 space-y-2 text-slate-400">
                <p>{t.detail}</p>
                {t.example && (
                  <p className="rounded bg-slate-900/60 p-2">
                    <span className="text-slate-500">ตัวอย่าง:</span> {t.example}
                  </p>
                )}
                {t.engineField && (
                  <p className="text-xs text-slate-500">
                    🔗 เชื่อมโยงในเอนจิน:{' '}
                    <code className="text-slate-400">{t.engineField}</code>
                  </p>
                )}
              </div>
            </details>
          </article>
        ))}
      </div>

      {filtered.length === 0 && (
        <p className="py-10 text-center text-slate-500">
          ไม่พบคำศัพท์ที่ตรงกับ “{query}”
        </p>
      )}
    </main>
  );
}
