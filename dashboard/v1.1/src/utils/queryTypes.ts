
export interface queryValueCPUUsage {
  CPUUsage: string;
  normalizedTime: number;
}

export interface queryValueMemoryUsedPercent {
  memoryUsedPercent: string;
  normalizedTime: number;
}

export interface chartData {
  x: number;
  y: string;
}

export interface QueryRange {
  start: number;
  end: number;
}

export interface queryValueDisk {
  cached: string;
  diskIO: string;
  normalizedTime: number;
}

export interface queryValueMemory {
  availableBytes: string;
  freeBytes: string;
  totalBytes: string;
  usedBytes: string;
  usedPercent: string;
  normalizedTime: number;
}

export interface queryValueSystemMetrics {
  cpuTotalUsage: string;
  disk: queryValueDisk;
  memory: queryValueMemory;
}

export interface QueryValues {
  normalizedTime: number;
  timestamp: string;
  value: queryValueSystemMetrics;
}

export interface QueryResponse {
  queryTime: number;
  range: QueryRange;
  timeSeriesPath: string;
  values: QueryValues[];
}