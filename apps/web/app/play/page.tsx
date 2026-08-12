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
  EventType,
  TileType,
  type FinancialStatement,
  type GameState,
  type Tile,
} from '@/lib/engine-wasm/types';

const PER_SIDE = 7; // 7×7 perimeter = 24 ช่อง

// แม็ป tile index → (row, col) ในกริด 7×7 เรียงตามเข็มนาฬิกาจากมุมบนซ้าย
// index 0 (Payday) = มุมบนซ้าย
function tileToGrid(i: number): { r: number; c: number } {
  if (i <= 6) return { r: 0, c: i }; // top L→R
  if (i <= 12) return { r: i - 6, c: 6 }; // right T→B
  if (i <= 18) return { r: 6, c: 18 - i }; // bottom R→L
  return { r: 24 - i, c: 0 }; // left B→T
}

const TILE_STYLE: Record<number, { icon: string; cls: string }> = {
  [TileType.Payday]: { icon: '💰', cls: 'border-emerald-500/60 bg-emerald-500/15' },
  [TileType.Opportunity]: { icon: '🃏', cls: 'border-sky-500/60 bg-sky-500/15' },
  [TileType.Shopping]: { icon: '🛍️', cls: 'border-amber-500/60 bg-amber-500/15' },
  [TileType.Crisis]: { icon: '⚠️', cls: 'border-red-500/60 bg-red-500/15' },
  [TileType.Market]: { icon: '📈', cls: 'border-violet-500/60 bg-violet-500/15' },
  [TileType.Downsizing]: { icon: '📉', cls: 'border-rose-500/60 bg-rose-500/15' },
  [TileType.Baby]: { icon: '👶', cls: 'border-pink-500/60 bg-pink-500/15' },
  [TileType.Charity]: { icon: '❤️', cls: 'border-red-400/60 bg-red-400/15' },
  [TileType.Blank]: { icon: '·', cls: 'border-slate-700 bg-slate-800/30' },
};

function tileIcon(type?: number): string {
  if (type === undefined) return '';
  return TILE_STYLE[type]?.icon ?? '·';
}

interface LogEntry {
  id: number;
  text: string;
}

const PLAYER_COLORS = ['bg-emerald-500', 'bg-sky-500', 'bg-violet-500', 'bg-amber-500'];

export default function PlayPage() {
  const [status, setStatus] = useState('กำลังโหลด engine…');
  const [board, setBoard] = useState<Tile[]>([]);
  const [state, setState] = useState<GameState | null>(null);
  const [statements, setStatements] = useState<FinancialStatement[]>([]);
  const [log, setLog] = useState<LogEntry[]>([]);
  const [rolling, setRolling] = useState(false);
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
    setLog((l) => [{ id: logId.current, text }, ...l].slice(0, 8));
  }

  // ดึง breakdown งบการเงินของทุกผู้เล่นจาก engine
  function fetchStatements(n: number): FinancialStatement[] {
    const out: FinancialStatement[] = [];
    for (let i = 0; i < n; i++) out.push(getStatement(i));
    return out;
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

  function handleRoll() {
    if (!state || rolling) return;
    setRolling(true);
    try {
      const cur = state.Players[state.CurrentTurn];
      const prevPos = cur.Position;
      const prevCash = cur.Cash;
      const n = board.length || 24;

      const events = applyAction({ PlayerID: cur.ID, Type: ActionType.Roll });
      const next = getGameState();
      const moved = next.Players.find((p) => p.ID === cur.ID);
      if (!moved) throw new Error('player not found after roll');

      const steps = ((moved.Position - prevPos) + n) % n;
      const gain = moved.Cash - prevCash;
      const gotPayday = events.some((e) => e.Type === EventType.Payday);

      let text = `🎲 ${cur.Name} ทอยได้ ${steps} → ช่อง ${moved.Position} ${tileIcon(board[moved.Position]?.Type)}`;
      if (gotPayday) text += `  →  💰 +${gain.toLocaleString()}`;
      pushLog(text);

      setState(next);
    } catch (e) {
      setStatus(`❌ ${(e as Error).message}`);
    } finally {
      setRolling(false);
    }
  }

  const current = state ? state.Players[state.CurrentTurn] : null;
  const currentStmt = state ? statements[state.CurrentTurn] : undefined;
  const ready = state !== null;

  return (
    <main className="mx-auto max-w-5xl px-6 py-10">
      <header className="mb-6">
        <h1 className="text-3xl font-bold text-emerald-400">🎮 เล่นเกม (Slice 1.5 + อาชีพจริง)</h1>
        <p className="mt-1 text-sm text-slate-400">
          ผู้เล่นถูกสุ่มอาชีพจริง (เงินเดือน/ภาษี/ประกันสังคม/หนี้) — กดทอยเต๋าเดินรอบกระดาน ผ่าน Payday แล้วรับเงินสดสุทธิ
        </p>
      </header>

      <div className="grid gap-6 lg:grid-cols-[minmax(0,1fr)_340px]">
        {/* ── กระดานจตุรัส ── */}
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

            {/* center panel */}
            <div
              style={{ gridRow: '2 / span 5', gridColumn: '2 / span 5' }}
              className="flex flex-col items-center justify-center gap-2 rounded-lg border border-slate-700 bg-slate-900/70 p-4 text-center"
            >
              <div className="text-sm font-semibold text-emerald-400">Finance Boardgame</div>
              {current ? (
                <div className="text-lg text-slate-100">
                  ถึงตา <span className="font-bold">{current.Name}</span>
                </div>
              ) : (
                <div className="text-slate-500">…</div>
              )}
              {currentStmt && (
                <div className="text-xs text-slate-400">
                  สุทธิ/เดือน{' '}
                  <b className="text-emerald-400">{currentStmt.MonthlyCashFlow.toLocaleString()}</b>{' '}
                  · รอบที่ {state?.Round ?? 0}
                </div>
              )}
              <button
                onClick={handleRoll}
                disabled={!ready || rolling}
                className="mt-1 rounded-md bg-emerald-600 px-5 py-2 font-medium text-white transition hover:bg-emerald-500 disabled:cursor-not-allowed disabled:opacity-40"
              >
                {rolling ? 'กำลังทอย…' : '🎲 ทอยเต๋า'}
              </button>
              <button
                onClick={startNewGame}
                disabled={rolling}
                className="text-xs text-slate-400 underline-offset-2 transition hover:text-emerald-400 hover:underline disabled:opacity-40"
              >
                สุ่มอาชีพใหม่
              </button>
              <p className="mt-1 max-w-xs text-xs text-slate-500">{status}</p>
            </div>

            {/* ── tokens ── */}
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

        {/* ── ข้าง: งบการเงินผู้เล่น + event log ── */}
        <aside className="flex flex-col gap-4">
          <section className="rounded-lg border border-slate-700 bg-slate-800/40 p-4">
            <h2 className="mb-3 text-sm font-semibold text-slate-300">ผู้เล่น & งบการเงิน</h2>
            <ul className="space-y-3">
              {state?.Players.map((p, idx) => {
                const fs = statements[idx];
                const isCurrent = current?.ID === p.ID;
                return (
                  <li
                    key={p.ID}
                    className={`rounded p-3 ${isCurrent ? 'bg-emerald-500/10 ring-1 ring-emerald-500/40' : 'bg-slate-900/40'}`}
                  >
                    <div className="flex items-center gap-2 text-sm">
                      <span className={`h-3 w-3 shrink-0 rounded-full ${PLAYER_COLORS[idx % PLAYER_COLORS.length]}`} />
                      <span className="font-semibold text-slate-100">{p.Name}</span>
                      <span className="ml-auto font-mono text-emerald-400">
                        {p.Cash.toLocaleString()}
                      </span>
                    </div>
                    {fs && (
                      <div className="mt-2 grid grid-cols-2 gap-x-3 gap-y-0.5 text-xs text-slate-400">
                        <span>เงินเดือน <b className="text-slate-300">{fs.EarnedIncome.toLocaleString()}</b></span>
                        <span>ภาษี <b className="text-slate-300">{fs.Tax.toLocaleString()}</b></span>
                        <span>ประกันสังคม <b className="text-slate-300">{fs.SocialSecurity.toLocaleString()}</b></span>
                        <span>รายจ่ายรวม <b className="text-slate-300">{fs.TotalExpenses.toLocaleString()}</b></span>
                        <span className="col-span-2">
                          สุทธิ/เดือน <b className="text-emerald-400">{fs.MonthlyCashFlow.toLocaleString()}</b>
                          <span className="ml-2 text-slate-500">· ช่อง {p.Position}</span>
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
                  <li key={e.id} className="border-l-2 border-slate-600 pl-2">
                    {e.text}
                  </li>
                ))}
              </ul>
            )}
          </section>
        </aside>
      </div>
    </main>
  );
}
