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
  value: queryValueSystemMetrics & ping & jitter & monitor & journal;
}

export interface QueryResponse {
  queryTime: number;
  range: QueryRange;
  timeSeriesPath: string;
  values: QueryValues[];
}

export interface journal {
  cerr: number;
  cevents: number;
  ckerr: number;
  ckevents: number;
  ckwarn: number;
  cwarn: number;
}

export interface APIQueryResponse {
  data: QueryResponse;
  success: string;
}

export interface PathData {
  matrixName: string;
  ping: string;
  jitter: string;
  fping: string;
  monitor: string;
}
export interface TimeSeries {
  name: string;
  path: PathData;
}

export interface APITimeSeriesResponse {
  data: TimeSeries[];
  success: string;
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

export interface Path {
  fping: string;
  jitter: string;
  monitor: string;
  ping: string;
  matrixName: string;
}

export interface TestServicesRoutes {
  testServicesRoutes: string[];
}

export interface RoutesSummary {
  testServicesRoutes: string[];
  monitoringRoutes: string[];
}

export interface TimeSeriesPath {
  name: string;
  path: Path;
}

export interface MatrixResponse {
  jitter: QueryResponse;
  monitor: QueryResponse;
  ping: QueryResponse;
}

export interface RouteDetails {
  ping: QueryResponse;
  jitter: QueryResponse;
  monitor: QueryResponse;
  name: string;
}
