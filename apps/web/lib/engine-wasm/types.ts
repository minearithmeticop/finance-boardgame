// ⚠️ TECH DEBT — จดไว้ใน NOTE.md (Slice 1.5)
// TS types นี้ **mirror Go domain ด้วยมือ** จาก packages/engine/domain/types.go
// เนื่องจาก Go ↔ TS ไม่ share type อัตโนมัติ และ enum ของ Go serialize เป็น int (iota)
// → มีความเสี่ยง drift (NFR-006)
// Slice ถัดไปจะใช้ codegen สร้างไฟล์นี้จาก Go โดยอัตโนมัติ แทนการเขียนมือ

// หมายเหตุ: field names เป็น PascalCase ตรง Go (Go ยังไม่มี json tag)

export interface GameState {
  Phase: number;
  Players: Player[];
  CurrentTurn: number;
  Round: number;
  Seed: number;
}

export interface Player {
  ID: string;
  Name: string;
  IsAI: boolean;
  Cash: number;
  Profession: Profession;
  Assets: Asset[];
  Liabilities: Liability[];
  Position: number;
  OnFastTrack: boolean;
  Bankrupt: boolean;
}

export interface Profession {
  Name: string;
  Salary: number;
  Taxes: number;
  SocialSecurity: number;
  OtherExpenses: number;
  HomeMortgage: Liability;
  SchoolLoan: Liability;
  CarLoan: Liability;
  CreditCard: Liability;
  Savings: number;
}

export interface Liability {
  ID: string;
  Name: string;
  Payment: number;
  Balance: number;
}

export interface Asset {
  ID: string;
  Type: number;
  Name: string;
  CashFlow: number;
  Cost: number;
  DownPayment: number;
  LoanRemaining: number;
}

export interface Event {
  Type: number;
  PlayerID: string;
}

export interface Tile {
  Type: number;
  Name: string;
}

// ── Enum const mirrors (ค่าตรงกับ iota ใน Go domain) ──────────────────────
// อ้างอิง: packages/engine/domain/types.go

export const Phase = {
  RatRace: 0,
  FastTrack: 1,
  Ended: 2,
} as const;

export const ActionType = {
  Roll: 0,
  BuyAsset: 1,
  SellAsset: 2,
  PayOffLiability: 3,
  TakeLoan: 4,
  EndTurn: 5,
} as const;

export const EventType = {
  Moved: 0,
  Landed: 1,
  CashChanged: 2,
  AssetBought: 3,
  AssetSold: 4,
  Payday: 5,
  EscapedRatRace: 6,
  GameWon: 7,
} as const;

export const TileType = {
  Payday: 0,
  Opportunity: 1,
  Doodad: 2,
  Market: 3,
  Downsizing: 4,
  Baby: 5,
  Charity: 6,
  Blank: 7,
} as const;
