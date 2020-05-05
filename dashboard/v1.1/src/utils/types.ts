export const HOST_IP = 'http://localhost:9090';

export interface service_states {
  ping: string;
  floodping: string;
  jitter: string;
  monitoring: string;
}

// System Value details
export interface SystemValue {
  range: rangeValues;
  queryTime: number;
  values: dataValues[];
}

interface rangeValues {
  start: any;
  end: any;
}

interface dataValues {
  Value: dataValue;
  timestamp: string;
  normalizedTimes: string;
}

interface dataValue {
  cpuTotalUsage: any;
  memory: memoryDetails;
  disk: diskDetails;
}

interface memoryDetails {
  totalBytes: string;
  availableBytes: string;
  usedBytes: string;
  usedPercent: any;
  freeBytes: any;
}

interface diskDetails {
  diskIO: string;
  cached: string;
}
