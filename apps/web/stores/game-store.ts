import { create } from 'zustand';

// Skeleton game store — จะ extend ใน Session ถัดไปเพื่อเก็บ game state
// ที่ sync ระหว่าง engine host (WASM หรือ backend) กับ UI
interface GameStoreState {
  /** version ของ engine ที่กำลังใช้งาน (WASM หรือ backend) */
  engineVersion: string;
  setEngineVersion: (v: string) => void;

  /** โหมดการเล่นปัจจุบัน */
  mode: 'online' | 'ai' | 'local' | 'async' | null;
  setMode: (m: GameStoreState['mode']) => void;
}

export const useGameStore = create<GameStoreState>((set) => ({
  engineVersion: 'unknown',
  setEngineVersion: (engineVersion) => set({ engineVersion }),

  mode: null,
  setMode: (mode) => set({ mode }),
}));
