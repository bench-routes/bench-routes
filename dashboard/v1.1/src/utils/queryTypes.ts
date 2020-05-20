
export interface queryValueCPUUsage {
  CPUUsage: string;
  normalizedTime: number;
}

export interface queryValueMemoryUsedPercent {
  memoryUsedPercent: string;
  normalizedTime: number;
}

export interface chartData {
  x: string | number;
  y: string | number;
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
  value: queryValueSystemMetrics & ping & jitter & monitor;
}

export interface QueryResponse {
  queryTime: number;
  range: QueryRange;
  timeSeriesPath: string;
  values: QueryValues[];
}

export interface ping {
  avgValue: string;
  maxValue: string;
  mdevValue: string;
  minValue: string;
}

export interface jitter {
  value: string;
}

export interface monitor {
  delay: number;
  resLength: number;
  resStatusCode: number;
}