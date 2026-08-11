// Client สำหรับสื่อสารกับ Go backend (ใช้ใน Online / Async mode)
// สำหรับ Local / AI mode ให้ใช้ lib/engine-wasm แทน (รัน engine ใน browser ตรงๆ)

const API_BASE = process.env.NEXT_PUBLIC_API_BASE ?? 'http://localhost:8080';

export interface HealthResponse {
  status: string;
  service: string;
  engine: string;
}

/** เรียก /healthz เพื่อเช็คว่า backend พร้อมใช้งาน */
export async function health(): Promise<HealthResponse> {
  const res = await fetch(`${API_BASE}/healthz`);
  if (!res.ok) throw new Error(`healthz failed: ${res.status}`);
  return res.json() as Promise<HealthResponse>;
}

// TODO(Session#3):
//   createGame(players), joinGame(roomId), getGameState(roomId)
//   sendAction(roomId, action) — async mode
//   connectWS(roomId) — real-time mode
