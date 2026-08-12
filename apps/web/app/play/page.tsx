'use client';

import { useEffect, useRef, useState } from 'react';
import {
  applyAction,
  createGameWithRandomProfessions,
  engineVersion,
  getBoard,
  getGameState,
  getStatement,
  loadEngineWasm,
} from '@/lib/engine-wasm/wasm';
import {
  ActionType,
  AssetType,
  EventType,
  TileType,
  type Event,
  type FinancialStatement,
  type GameState,
  type Tile,
} from '@/lib/engine-wasm/types';

const PER_SIDE = 7; // 7×7 perimeter = 24 ช่อง

function tileToGrid(i: number): { r: number; c: number } {
  if (i <= 6) return { r: 0, c: i };
  if (i <= 12) return { r: i - 6, c: 6 };
  if (i <= 18) return { r: 6, c: 18 - i };
  return { r: 24 - i, c: 0 };
}

const TILE_STYLE: Record<number, { icon: string; cls: string }> = {
  [TileType.Payday]: { icon: '💰', cls: 'border-emerald-500/60 bg-emerald-500/15' },
  [TileType.Opportunity]: { icon: '🃏', cls: 'border-sky-500/60 bg-sky-500/15' },
  [TileType.Shopping]: { icon: '🛍️', cls: 'border-amber-500/60 bg-amber-500/15' },
  [TileType.Crisis]: { icon: '⚠️', cls: 'border-red-500/60 bg-red-500/15' },
  [TileType.Market]: { icon: '📈', cls: 'border-violet-500/60 bg-violet-500/15' },
  [TileType.Downsizing]: { icon: '📉', cls: 'border-rose-500/60 bg-rose-500/15' },
  [TileType.Family]: { icon: '👨‍👩‍👧', cls: 'border-pink-500/60 bg-pink-500/15' },
  [TileType.Donate]: { icon: '❤️', cls: 'border-rose-400/60 bg-rose-400/15' },
  [TileType.Blank]: { icon: '·', cls: 'border-slate-700 bg-slate-800/30' },
  [TileType.News]: { icon: '📰', cls: 'border-slate-400/60 bg-slate-500/15' },
  [TileType.Windfall]: { icon: '🎁', cls: 'border-teal-500/60 bg-teal-500/15' },
  [TileType.SideJob]: { icon: '💼', cls: 'border-indigo-500/60 bg-indigo-500/15' },
  [TileType.Learn]: { icon: '📚', cls: 'border-cyan-500/60 bg-cyan-500/15' },
  [TileType.Health]: { icon: '🩺', cls: 'border-orange-500/60 bg-orange-500/15' },
};

// หมวด LifeEvent → ไอคอน/ป้าย สำหรับ event log
const CATEGORY_LOG: Record<string, string> = {
  news: '📰', windfall: '🎁', sidejob: '💼', shopping: '🛍️',
  family: '👨‍👩‍👧', donate: '❤️', learn: '📚', health: '🩺', crisis: '⚠️',
};

const ASSET_LABEL: Record<number, string> = {
  [AssetType.Stock]: 'หุ้น',
  [AssetType.RealEstate]: 'อสังหาฯ',
  [AssetType.Business]: 'ธุรกิจ',
  [AssetType.Other]: 'อื่นๆ',
};

const PLAYER_COLORS = ['bg-emerald-500', 'bg-sky-500', 'bg-violet-500', 'bg-amber-500'];

function tileIcon(type?: number): string {
  if (type === undefined) return '';
  return TILE_STYLE[type]?.icon ?? '·';
}

interface LogEntry {
  id: number;
  text: string;
}

// สร้างบรรทัด log จาก events ของ action หนึ่ง
function logFromEvents(events: Event[], name: string, board: Tile[]): string {
  const parts: string[] = [];
  let moved = '';
  for (const ev of events) {
    const d = (ev.Data ?? {}) as Record<string, number | string>;
    switch (ev.Type) {
      case EventType.Moved: {
        const pos = Number(d.position);
        moved = `🎲 ${name} ทอยได้ ${d.steps} → ช่อง ${pos} ${tileIcon(board[pos]?.Type)}`;
        break;
      }
      case EventType.Payday:
        parts.push(`💰 +${Number(d.amount).toLocaleString()}`);
        break;
      case EventType.Landed:
        if (d.kind === 'opportunity') parts.push(`🃏 ดีล: ${d.title}`);
        else if (d.kind === 'declined') parts.push(`⏭️ ผ่าน: ${d.title}`);
        break;
      case EventType.CashChanged: {
        const kind = String(d.kind ?? '');
        const icon = CATEGORY_LOG[kind];
        if (icon) {
          const amt = Number(d.amount ?? 0);
          parts.push(amt === 0 ? `${icon} ${d.title}` : `${icon} ${d.title} ${amt.toLocaleString()}`);
        }
        break;
      }
      case EventType.AssetBought:
        parts.push(`🏠 ซื้อ: ${d.title}`);
        break;
    }
  }
  return [moved, ...parts].filter(Boolean).join('   ');
}

export default function PlayPage() {
  const [status, setStatus] = useState('กำลังโหลด engine…');
  const [board, setBoard] = useState<Tile[]>([]);
  const [state, setState] = useState<GameState | null>(null);
  const [statements, setStatements] = useState<FinancialStatement[]>([]);
  const [log, setLog] = useState<LogEntry[]>([]);
  const [busy, setBusy] = useState(false);
  const [loanOpen, setLoanOpen] = useState(false);
  const [loanType, setLoanType] = useState('personal');
  const [loanAmount, setLoanAmount] = useState('');
  const [collatSel, setCollatSel] = useState('home');
  const [loanMsg, setLoanMsg] = useState('');
  const logId = useRef(0);

  useEffect(() => {
    let cancelled = false;
    (async () => {
      try {
        await loadEngineWasm();
        if (cancelled) return;
        setBoard(getBoard());
        startNewGame();
        setStatus(`✅ engine v${engineVersion()} พร้อม`);
      } catch (e) {
        setStatus(`❌ ${(e as Error).message} — ลองรัน \`pnpm build:wasm\` ก่อน`);
      }
    })();
    return () => {
      cancelled = true;
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  function pushLog(text: string) {
    logId.current += 1;
    setLog((l) => [{ id: logId.current, text }, ...l].slice(0, 12));
  }

  function fetchStatements(n: number): FinancialStatement[] {
    const out: FinancialStatement[] = [];
    for (let i = 0; i < n; i++) out.push(getStatement(i));
    return out;
  }

  function refreshFrom(state: GameState, refreshStmts: boolean) {
    setState(state);
    if (refreshStmts) setStatements(fetchStatements(state.Players.length));
  }

  function startNewGame() {
    try {
      const seed = Math.floor(Math.random() * 1_000_000);
      createGameWithRandomProfessions(seed, 2);
      const s = getGameState();
      setState(s);
      setStatements(fetchStatements(s.Players.length));
      setLog([]);
      setStatus(`✅ เริ่มเกมใหม่ (seed ${seed}) — ${s.Players.map((p) => p.Name).join(' vs ')}`);
    } catch (e) {
      setStatus(`❌ ${(e as Error).message}`);
    }
  }

  function nameOf(id: string): string {
    return state?.Players.find((p) => p.ID === id)?.Name ?? id;
  }

  function handleRoll() {
    if (!state || busy || state.Pending) return;
    setBusy(true);
    try {
      const cur = state.Players[state.CurrentTurn];
      const events = applyAction({ PlayerID: cur.ID, Type: ActionType.Roll });
      const next = getGameState();
      pushLog(logFromEvents(events, cur.Name, board));
      refreshFrom(next, false);
    } catch (e) {
      setStatus(`❌ ${(e as Error).message}`);
    } finally {
      setBusy(false);
    }
  }

  function resolvePending(type: number, verb: string, refreshStmts: boolean) {
    if (!state?.Pending) return;
    setBusy(true);
    try {
      const pending = state.Pending;
      applyAction({ PlayerID: pending.PlayerID, Type: type });
      const next = getGameState();
      pushLog(`${verb} ${nameOf(pending.PlayerID)}: ${pending.DealCard.Title}`);
      refreshFrom(next, refreshStmts);
    } catch (e) {
      setStatus(`❌ ${(e as Error).message}`);
    } finally {
      setBusy(false);
    }
  }

  const handleBuy = () => resolvePending(ActionType.BuyAsset, '🏠 ซื้อ', true);
  const handleDecline = () => resolvePending(ActionType.Decline, '⏭️ ผ่าน', false);

  function handleTakeLoan() {
    if (!state) return;
    const cur = state.Players[state.CurrentTurn];
    try {
      const payload: Record<string, unknown> = { lender: loanType, amount: Number(loanAmount) };
      if (loanType === 'secured') {
        if (collatSel === 'home' || collatSel === 'car') {
          payload.collateralKind = collatSel;
        } else {
          payload.collateralKind = 'asset';
          payload.collateralRef = collatSel;
        }
      }
      applyAction({ PlayerID: cur.ID, Type: ActionType.TakeLoan, Payload: payload });
      const next = getGameState();
      setState(next);
      setStatements(fetchStatements(next.Players.length));
      setLoanMsg(`✅ อนุมัติ! ได้รับ ${Number(loanAmount).toLocaleString()} บาท`);
      setLoanAmount('');
    } catch (e) {
      setLoanMsg(`❌ ${(e as Error).message}`);
    }
  }

  function handlePayOffLoan(loanID: string) {
    if (!state) return;
    const cur = state.Players[state.CurrentTurn];
    try {
      applyAction({ PlayerID: cur.ID, Type: ActionType.PayOffLiability, Payload: { loanID } });
      const next = getGameState();
      setState(next);
      setStatements(fetchStatements(next.Players.length));
      setLoanMsg('✅ ปิดสินเชื่อแล้ว');
    } catch (e) {
      setLoanMsg(`❌ ${(e as Error).message}`);
    }
  }

  const current = state ? state.Players[state.CurrentTurn] : null;
  const currentStmt = state ? statements[state.CurrentTurn] : undefined;
  const pending = state?.Pending ?? null;
  const pendingPlayer = pending ? state!.Players.find((p) => p.ID === pending.PlayerID) : undefined;
  const canAfford = pending && pendingPlayer
    ? pendingPlayer.Cash >= pending.DealCard.DownPayment
    : false;
  const pendingNet = pending ? pending.DealCard.CashFlow - pending.DealCard.LoanPayment : 0;

  return (
    <main className="mx-auto max-w-5xl px-6 py-10">
      <header className="mb-6">
        <h1 className="text-3xl font-bold text-emerald-400">🎮 เล่นเกม (Slice 3 — ดีล + วิกฤต)</h1>
        <p className="mt-1 text-sm text-slate-400">
          สุ่มอาชีพจริง → ทอยเต๋า → ตก 🃏 Opportunity (เลือกซื้อสินทรัพย์) / 🛍️ Shopping / ⚠️ Crisis
        </p>
      </header>

      <div className="grid gap-6 lg:grid-cols-[minmax(0,1fr)_340px]">
        {/* ── กระดาน ── */}
        <div className="relative mx-auto aspect-square w-full max-w-2xl">
          <div className="relative grid h-full w-full grid-cols-7 grid-rows-7 gap-1">
            {board.map((tile, i) => {
              const { r, c } = tileToGrid(i);
              const st = TILE_STYLE[tile.Type] ?? TILE_STYLE[TileType.Blank];
              return (
                <div
                  key={i}
                  style={{ gridRow: r + 1, gridColumn: c + 1 }}
                  className={`flex items-center justify-center rounded border text-base ${st.cls}`}
                  title={`ช่อง ${i}: ${tile.Name}`}
                >
                  <span>{st.icon}</span>
                </div>
              );
            })}

            <div
              style={{ gridRow: '2 / span 5', gridColumn: '2 / span 5' }}
              className="flex flex-col items-center justify-center gap-2 rounded-lg border border-slate-700 bg-slate-900/70 p-4 text-center"
            >
              <div className="text-sm font-semibold text-emerald-400">Finance Boardgame</div>
              {current ? (
                <div className="text-lg text-slate-100">
                  {pending ? '🃏 ตัดสินใจดีล' : 'ถึงตา'} <span className="font-bold">{current.Name}</span>
                </div>
              ) : (
                <div className="text-slate-500">…</div>
              )}
              {currentStmt && (
                <div className="text-xs text-slate-400">
                  สุทธิ/เดือน <b className="text-emerald-400">{currentStmt.MonthlyCashFlow.toLocaleString()}</b>
                  <span className="ml-2 text-slate-500">· รอบที่ {state?.Round ?? 0}</span>
                </div>
              )}
              {!pending && (
                <button
                  onClick={handleRoll}
                  disabled={busy}
                  className="mt-1 rounded-md bg-emerald-600 px-5 py-2 font-medium text-white transition hover:bg-emerald-500 disabled:cursor-not-allowed disabled:opacity-40"
                >
                  {busy ? '…' : '🎲 ทอยเต๋า'}
                </button>
              )}
              <button
                onClick={() => { setLoanOpen(true); setLoanMsg(''); }}
                disabled={busy}
                className="text-xs text-sky-300 transition hover:text-sky-200 disabled:opacity-40"
              >
                💳 สินเชื่อ
              </button>
              <button
                onClick={startNewGame}
                disabled={busy}
                className="text-xs text-slate-400 underline-offset-2 transition hover:text-emerald-400 hover:underline disabled:opacity-40"
              >
                สุ่มอาชีพใหม่
              </button>
              <p className="mt-1 max-w-xs text-xs text-slate-500">{status}</p>
            </div>

            {state?.Players.map((p, idx) => {
              const { r, c } = tileToGrid(p.Position);
              const onTile = state.Players.filter((q) => q.Position === p.Position);
              const offset = onTile.length > 1 ? (onTile.indexOf(p) === 0 ? -1.2 : 1.2) : 0;
              const left = ((c + 0.5) / PER_SIDE) * 100 + offset;
              const top = ((r + 0.5) / PER_SIDE) * 100;
              const isCurrent = current?.ID === p.ID;
              return (
                <div
                  key={p.ID}
                  style={{
                    left: `${left}%`,
                    top: `${top}%`,
                    transform: 'translate(-50%, -50%)',
                    transition: 'left 0.4s ease, top 0.4s ease',
                  }}
                  className={`pointer-events-none absolute z-10 flex h-6 w-6 items-center justify-center rounded-full border-2 border-white/80 text-[10px] font-bold text-white shadow-lg ${PLAYER_COLORS[idx % PLAYER_COLORS.length]} ${isCurrent ? 'scale-125 ring-2 ring-emerald-300' : ''}`}
                  title={p.Name}
                >
                  {idx + 1}
                </div>
              );
            })}
          </div>
        </div>

        {/* ── ข้าง ── */}
        <aside className="flex flex-col gap-4">
          <section className="rounded-lg border border-slate-700 bg-slate-800/40 p-4">
            <h2 className="mb-3 text-sm font-semibold text-slate-300">ผู้เล่น & งบการเงิน</h2>
            <ul className="space-y-3">
              {state?.Players.map((p, idx) => {
                const fs = statements[idx];
                const isCurrent = current?.ID === p.ID;
                const assets = p.Assets ?? [];
                return (
                  <li
                    key={p.ID}
                    className={`rounded p-3 ${isCurrent ? 'bg-emerald-500/10 ring-1 ring-emerald-500/40' : 'bg-slate-900/40'}`}
                  >
                    <div className="flex items-center gap-2 text-sm">
                      <span className={`h-3 w-3 shrink-0 rounded-full ${PLAYER_COLORS[idx % PLAYER_COLORS.length]}`} />
                      <span className="font-semibold text-slate-100">{p.Name}</span>
                      <span className="ml-auto font-mono text-emerald-400">{p.Cash.toLocaleString()}</span>
                    </div>
                    {fs && (
                      <div className="mt-2 grid grid-cols-2 gap-x-3 gap-y-0.5 text-xs text-slate-400">
                        <span>เงินเดือน <b className="text-slate-300">{fs.EarnedIncome.toLocaleString()}</b></span>
                        <span>สุทธิ/เดือน <b className="text-emerald-400">{fs.MonthlyCashFlow.toLocaleString()}</b></span>
                        <span className="col-span-2 text-slate-500">
                          สินทรัพย์ {assets.length} หน่วย · ช่อง {p.Position}
                          {assets.length > 0 && (
                            <span className="ml-1 text-slate-400">({assets.map((a) => a.Name).join(', ')})</span>
                          )}
                        </span>
                      </div>
                    )}
                  </li>
                );
              })}
            </ul>
          </section>

          <section className="rounded-lg border border-slate-700 bg-slate-800/40 p-4">
            <h2 className="mb-2 text-sm font-semibold text-slate-300">ล่าสุด</h2>
            {log.length === 0 ? (
              <p className="text-xs text-slate-500">ยังไม่มีการทอย — กด “🎲 ทอยเต๋า”</p>
            ) : (
              <ul className="space-y-1 text-xs text-slate-300">
                {log.map((e) => (
                  <li key={e.id} className="border-l-2 border-slate-600 pl-2">{e.text}</li>
                ))}
              </ul>
            )}
          </section>
        </aside>
      </div>

      {/* ── Loan modal ── */}
      {loanOpen && state && (() => {
        const cur = state.Players[state.CurrentTurn];
        const collatOpts: { sel: string; label: string; value: number }[] = [];
        if (cur.Profession.HomeMortgage.Balance > 0)
          collatOpts.push({ sel: 'home', label: `ค้ำบ้าน (${cur.Profession.HomeMortgage.Balance.toLocaleString()})`, value: cur.Profession.HomeMortgage.Balance });
        if (cur.Profession.CarLoan.Balance > 0)
          collatOpts.push({ sel: 'car', label: `ค้ำรถ (${cur.Profession.CarLoan.Balance.toLocaleString()})`, value: cur.Profession.CarLoan.Balance });
        cur.Assets.forEach((a) => collatOpts.push({ sel: a.ID, label: `ค้ำ ${a.Name} (${a.Cost.toLocaleString()})`, value: a.Cost }));
        const selOpt = collatOpts.find((o) => o.sel === collatSel) ?? collatOpts[0];
        const maxLoan = selOpt ? Math.floor(selOpt.value * 0.7) : 0;
        return (
          <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 p-4">
            <div className="max-h-[90vh] w-full max-w-md overflow-y-auto rounded-lg border border-sky-500/50 bg-slate-900 p-5 shadow-xl">
              <div className="mb-1 flex items-center justify-between">
                <h3 className="text-lg font-bold text-slate-100">💳 สินเชื่อ — {cur.Name}</h3>
                <button onClick={() => setLoanOpen(false)} className="text-slate-500 hover:text-slate-300">✕</button>
              </div>
              <p className="mb-3 text-xs text-slate-400">เงินสด {cur.Cash.toLocaleString()} · สินเชื่อ {cur.Loans.length} บัญชี</p>
              <div className="mb-3 grid grid-cols-3 gap-2 text-center text-xs">
                {[
                  { id: 'personal', l: 'ส่วนบุคคล', r: '24%/ปี' },
                  { id: 'secured', l: 'ค้ำหลัก', r: '7%/ปี' },
                  { id: 'informal', l: 'นอกระบบ', r: '10%/ด!' },
                ].map((t) => (
                  <button
                    key={t.id}
                    onClick={() => { setLoanType(t.id); setLoanMsg(''); }}
                    className={`rounded border p-2 ${loanType === t.id ? 'border-sky-500 bg-sky-500/20 text-sky-300' : 'border-slate-700 text-slate-400'}`}
                  >
                    <div className="font-semibold">{t.l}</div>
                    <div className="text-[10px]">{t.r}</div>
                  </button>
                ))}
              </div>
              {loanType === 'informal' && (
                <p className="mb-2 rounded bg-red-500/15 p-2 text-xs text-red-300">⚠️ ดอกเบี้ย 10%/เดือน (120%/ปี) — กับดักหนี้สิน ใช้เป็นทางสุดท้าย!</p>
              )}
              {loanType === 'secured' && (
                <div className="mb-2">
                  {collatOpts.length === 0 ? (
                    <p className="text-xs text-rose-400">ไม่มีหลักค้ำ (ต้องมีบ้าน/รถ/สินทรัพย์)</p>
                  ) : (
                    <>
                      <label className="text-xs text-slate-400">หลักค้ำ (สูงสุด {maxLoan.toLocaleString()}):</label>
                      <select
                        value={collatSel}
                        onChange={(e) => setCollatSel(e.target.value)}
                        className="mt-1 w-full rounded border border-slate-700 bg-slate-800 px-2 py-1 text-sm text-slate-200"
                      >
                        {collatOpts.map((o) => (
                          <option key={o.sel} value={o.sel}>{o.label}</option>
                        ))}
                      </select>
                    </>
                  )}
                </div>
              )}
              <div className="mb-2 flex gap-2">
                <input
                  type="number"
                  value={loanAmount}
                  onChange={(e) => setLoanAmount(e.target.value)}
                  placeholder="วงเงินที่ต้องการ"
                  className="flex-1 rounded border border-slate-700 bg-slate-800 px-3 py-2 text-sm text-slate-100"
                />
                <button
                  onClick={handleTakeLoan}
                  disabled={busy || !loanAmount}
                  className="rounded bg-sky-600 px-4 py-2 text-sm font-medium text-white hover:bg-sky-500 disabled:opacity-40"
                >
                  ขอสินเชื่อ
                </button>
              </div>
              {loanMsg && <p className="mb-2 text-xs text-slate-300">{loanMsg}</p>}
              {cur.Loans.length > 0 && (
                <div className="mt-3 border-t border-slate-700 pt-3">
                  <div className="mb-2 text-xs font-semibold text-slate-400">สินเชื่อที่มี</div>
                  <ul className="space-y-2">
                    {cur.Loans.map((ln) => (
                      <li key={ln.ID} className="flex items-center justify-between rounded bg-slate-800/60 p-2 text-xs">
                        <span>
                          <b className="text-slate-200">{ln.Lender}</b>
                          {ln.Collateral && <span className="text-slate-500"> ({ln.Collateral})</span>}
                          <span className="ml-2 text-slate-400">ค่างวด {ln.MonthlyPay.toLocaleString()}/ด · คงค้าง {ln.Balance.toLocaleString()}</span>
                        </span>
                        <button
                          onClick={() => handlePayOffLoan(ln.ID)}
                          disabled={cur.Cash < ln.Balance || busy}
                          className="rounded bg-rose-600 px-2 py-1 text-white hover:bg-rose-500 disabled:opacity-40"
                        >
                          ปิด {ln.Balance.toLocaleString()}
                        </button>
                      </li>
                    ))}
                  </ul>
                </div>
              )}
            </div>
          </div>
        );
      })()}

      {/* ── Opportunity decision modal ── */}
      {pending && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 p-4">
          <div className="w-full max-w-sm rounded-lg border border-sky-500/50 bg-slate-900 p-5 shadow-xl">
            <div className="mb-1 text-xs text-sky-400">🃏 โอกาส — {ASSET_LABEL[pending.DealCard.AssetType] ?? 'สินทรัพย์'}</div>
            <h3 className="mb-3 text-xl font-bold text-slate-100">{pending.DealCard.Title}</h3>
            <dl className="grid grid-cols-2 gap-x-3 gap-y-1 text-sm text-slate-300">
              <dt className="text-slate-500">ราคาเต็ม</dt><dd className="text-right">{pending.DealCard.Cost.toLocaleString()}</dd>
              <dt className="text-slate-500">เงินดาวน์</dt><dd className="text-right">{pending.DealCard.DownPayment.toLocaleString()}</dd>
              <dt className="text-slate-500">รายได้/เดือน</dt><dd className="text-right text-emerald-400">+{pending.DealCard.CashFlow.toLocaleString()}</dd>
              {pending.DealCard.LoanPayment > 0 && (
                <><dt className="text-slate-500">ผ่อน/เดือน</dt><dd className="text-right text-rose-400">−{pending.DealCard.LoanPayment.toLocaleString()}</dd></>
              )}
              <dt className="text-slate-500">สุทธิ/เดือน</dt>
              <dd className={`text-right font-bold ${pendingNet >= 0 ? 'text-emerald-400' : 'text-rose-400'}`}>
                {pendingNet >= 0 ? '+' : ''}{pendingNet.toLocaleString()}
              </dd>
            </dl>
            {pendingPlayer && (
              <p className="mt-3 text-xs text-slate-500">
                เงินสด {pendingPlayer.Name}: <b className="text-slate-300">{pendingPlayer.Cash.toLocaleString()}</b>
                {!canAfford && <span className="ml-1 text-rose-400">(ไม่พอจ่ายดาวน์)</span>}
              </p>
            )}
            <div className="mt-4 flex gap-2">
              <button
                onClick={handleBuy}
                disabled={!canAfford || busy}
                className="flex-1 rounded-md bg-emerald-600 px-4 py-2 font-medium text-white transition hover:bg-emerald-500 disabled:cursor-not-allowed disabled:opacity-40"
              >
                ซื้อ (−{pending.DealCard.DownPayment.toLocaleString()})
              </button>
              <button
                onClick={handleDecline}
                disabled={busy}
                className="rounded-md border border-slate-600 px-4 py-2 font-medium text-slate-300 transition hover:bg-slate-800 disabled:opacity-40"
              >
                ผ่าน
              </button>
            </div>
          </div>
        </div>
      )}
    </main>
  );
}
