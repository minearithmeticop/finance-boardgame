// โหลด Game Engine (Go) ในรูปแบบ WebAssembly เพื่อรันฝั่ง browser
// สำหรับ Local Pass-and-Play และ Single Player vs AI (offline ได้)
//
// ไฟล์ WASM (engine.wasm + wasm_exec.js) generate โดย `pnpm build:wasm`
// (tooling/build-wasm.sh) แล้ววางใน /public/wasm — ถ้ายังไม่ build จะโหลดไม่สำเร็จ

declare global {
  interface Window {
    // คลาส Go runtime support ที่ wasm_exec.js ฉีดเข้ามา
    Go?: new () => {
      importObject: WebAssembly.Imports;
      run(instance: WebAssembly.Instance): Promise<void>;
    };
    // ฟังก์ชันที่ engine ลงทะเบียนไว้ใน cmd/wasm/main.go
    engineVersion?: () => string;
  }
}

let loadPromise: Promise<void> | null = null;

/** โหลด engine WASM ครั้งเดียว (memoized) แล้ว engineVersion() จะใช้ได้ */
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
    // go.run ลงทะเบียน globalThis.engineVersion แล้ว block (รัน main ของ WASM)
    go.run(result.instance);
  })().catch((err) => {
    // ให้ลองใหม่ได้ครั้งถัดไปถ้า fail
    loadPromise = null;
    throw err;
  });
  return loadPromise;
}

/** คืน version ของ engine (หลัง loadEngineWasm สำเร็จ) */
export function engineVersion(): string {
  return window.engineVersion?.() ?? 'wasm-not-loaded';
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
