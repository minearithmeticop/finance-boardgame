// โหลด Game Engine (Go) ในรูปแบบ WebAssembly + wrapper functions
// สำหรับ Local Pass-and-Play และ Single Player vs AI (offline ได้)
//
// ไฟล์ WASM (engine.wasm + wasm_exec.js) generate โดย `pnpm build:wasm`
// (tooling/build-wasm.sh) แล้ววางใน /public/wasm

import type { Event, GameState, Tile } from './types';

declare global {
  interface Window {
    // คลาส Go runtime support ที่ wasm_exec.js ฉีดเข้ามา
    Go?: new () => {
      importObject: WebAssembly.Imports;
      run(instance: WebAssembly.Instance): Promise<void>;
    };
    // engine API ที่ cmd/wasm/main.go ลงทะเบียนไว้
    engineVersion?: () => string;
    engineCreate?: (seed: number, playersJSON: string) => string;
    engineState?: () => string;
    engineApply?: (actionJSON: string) => string;
    engineBoard?: () => string;
    // flag "พร้อม" ที่ Go main ตั้งหลัง register ทุก callback
    __engineWasmReady?: boolean;
  }
}

let loadPromise: Promise<void> | null = null;

/** โหลด engine WASM ครั้งเดียว (memoized) แล้ว engine API จะใช้ได้ */
export function loadEngineWasm(): Promise<void> {
  if (loadPromise) return loadPromise;
  loadPromise = (async () => {
    // 1. โหลด wasm_exec.js (Go runtime support) ถ้ายังไม่มี
    if (!window.Go) {
      await loadScript('/wasm/wasm_exec.js');
    }
    const Go = window.Go;
    if (!Go) throw new Error('wasm_exec.js did not define global Go');

    // 2. fetch + instantiate engine.wasm
    const go = new Go();
    const resp = await fetch('/wasm/engine.wasm');
    if (!resp.ok) throw new Error(`fetch engine.wasm failed: ${resp.status}`);
    const bytes = new Uint8Array(await resp.arrayBuffer());
    const result = await WebAssembly.instantiate(bytes, go.importObject);
    // go.run ลงทะเบียน globalThis.engine* + __engineWasmReady แล้ว block
    go.run(result.instance);

    // 3. รอจนกว่า Go main จะ register callback ครบ (กัน race ตอน startup)
    await waitForReady();
  })().catch((err) => {
    loadPromise = null; // ให้ลองใหม่ได้ครั้งถัดไปถ้า fail
    throw err;
  });
  return loadPromise;
}

/** poll จนกว่า Go จะตั้ง __engineWasmReady (หมดเวลา 5s) */
async function waitForReady(timeoutMs = 5000): Promise<void> {
  const start = Date.now();
  while (!window.__engineWasmReady) {
    if (Date.now() - start > timeoutMs) {
      throw new Error('engine WASM did not become ready in time');
    }
    await new Promise((r) => setTimeout(r, 30));
  }
}

/** คืน version ของ engine (หลัง loadEngineWasm สำเร็จ) */
export function engineVersion(): string {
  return window.engineVersion?.() ?? 'wasm-not-loaded';
}

// ── engine wrappers ──────────────────────────────────────────────────────
// ทุกฟังก์ชัน Go คืน JSON envelope { data, error } → แกะแล้ว throw ถ้ามี error

function call<T>(fn: ((...a: any[]) => string) | undefined, ...args: unknown[]): T {
  if (typeof fn !== 'function') {
    throw new Error('engine not loaded — call loadEngineWasm() first');
  }
  const res = JSON.parse(fn(...args)) as { data: T; error: string };
  if (res.error) throw new Error(res.error);
  return res.data;
}

/** สร้างเกมใหม่ด้วย seed + รายชื่อผู้เล่น */
export function createGame(seed: number, players: unknown[]): void {
  call<void>(window.engineCreate, seed, JSON.stringify(players));
}

/** ดึง snapshot สถานะเกมปัจจุบัน */
export function getGameState(): GameState {
  return call<GameState>(window.engineState);
}

/** ส่ง action เข้า engine แล้วคืน events ที่เกิดขึ้น */
export function applyAction(action: unknown): Event[] {
  return call<Event[]>(window.engineApply, JSON.stringify(action));
}

/** ดึง layout กระดาน Rat Race (จาก engine ตรงๆ — single source) */
export function getBoard(): Tile[] {
  return call<Tile[]>(window.engineBoard);
}

function loadScript(src: string): Promise<void> {
  return new Promise((resolve, reject) => {
    const s = document.createElement('script');
    s.src = src;
    s.async = true;
    s.onload = () => resolve();
    s.onerror = () => reject(new Error(`failed to load ${src}`));
    document.head.appendChild(s);
  });
}
