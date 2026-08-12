'use client';

import { useEffect, useRef, useState } from 'react';
import {
  applyAction,
  createGame,
  engineVersion,
  getBoard,
  getGameState,
  loadEngineWasm,
} from '@/lib/engine-wasm/wasm';
import {
  ActionType,
  EventType,
  TileType,
  type GameState,
  type Player,
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
  [TileType.Doodad]: { icon: '🛍️', cls: 'border-amber-500/60 bg-amber-500/15' },
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

function makePlayer(id: string, name: string): Player {
  const zeroLiability = { ID: '', Name: '', Payment: 0, Balance: 0 };
  return {
    ID: id,
    Name: name,
    IsAI: false,
    Cash: 1000,
    Profession: {
      Name: 'Tester',
      Salary: 2000,
      Taxes: 0,
      OtherExpenses: 0,
      HomeMortgage: zeroLiability,
      SchoolLoan: zeroLiability,
      CarLoan: zeroLiability,
      CreditCard: zeroLiability,
      Savings: 0,
    },
    Assets: [],
    Liabilities: [],
    Position: 0,
    OnFastTrack: false,
    Bankrupt: false,
  };
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

  function startNewGame() {
    try {
      const seed = Math.floor(Math.random() * 1_000_000);
      createGame(seed, [makePlayer('p1', 'ผู้เล่น 1'), makePlayer('p2', 'ผู้เล่น 2')]);
      setState(getGameState());
      setLog([]);
      setStatus(`✅ เริ่มเกมใหม่ (seed ${seed}) — กดทอยเต๋า`);
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
  const ready = state !== null;

  return (
    <main className="mx-auto max-w-5xl px-6 py-10">
      <header className="mb-6">
        <h1 className="text-3xl font-bold text-emerald-400">🎮 เล่นเกม (Slice 1.5)</h1>
        <p className="mt-1 text-sm text-slate-400">
          engine รันใน browser ผ่าน WASM — กดทอยเต๋าเพื่อเดินรอบกระดาน ผ่าน Payday แล้วรับเงินสด
        </p>
      </header>

      <div className="grid gap-6 lg:grid-cols-[minmax(0,1fr)_320px]">
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
              className="flex flex-col items-center justify-center gap-3 rounded-lg border border-slate-700 bg-slate-900/70 p-4 text-center"
            >
              <div className="text-sm font-semibold text-emerald-400">Finance Boardgame</div>
              {current ? (
                <div className="text-lg text-slate-100">
                  ถึงตา <span className="font-bold">{current.Name}</span>
                </div>
              ) : (
                <div className="text-slate-500">…</div>
              )}
              <div className="text-xs text-slate-400">รอบที่ {state?.Round ?? 0}</div>
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
                เริ่มเกมใหม่
              </button>
              <p className="mt-1 text-xs text-slate-500">{status}</p>
            </div>

            {/* ── tokens (absolute, animate ตอน position เปลี่ยน) ── */}
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

        {/* ── ข้าง: สรุปผู้เล่น + event log ── */}
        <aside className="flex flex-col gap-4">
          <section className="rounded-lg border border-slate-700 bg-slate-800/40 p-4">
            <h2 className="mb-2 text-sm font-semibold text-slate-300">ผู้เล่น</h2>
            <ul className="space-y-2">
              {state?.Players.map((p, idx) => {
                const isCurrent = current?.ID === p.ID;
                return (
                  <li
                    key={p.ID}
                    className={`flex items-center gap-2 rounded p-2 text-sm ${isCurrent ? 'bg-emerald-500/10 ring-1 ring-emerald-500/40' : ''}`}
                  >
                    <span className={`h-3 w-3 rounded-full ${PLAYER_COLORS[idx % PLAYER_COLORS.length]}`} />
                    <span className="flex-1 text-slate-200">{p.Name}</span>
                    <span className="text-slate-400">ช่อง {p.Position}</span>
                    <span className="font-mono text-emerald-400">{p.Cash.toLocaleString()}</span>
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
