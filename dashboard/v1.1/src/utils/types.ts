export const HOST_IP = 'http://localhost:9990';

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

export interface routeEntryType {
    route: string,
    options: routeOptionsInterface[] | undefined
}

export interface routeOptionsInterface {
  params: { Name: string; Value: string }[],
  headers: { OfType: string; Value: string }[],
  method: string,
  body: { Name: string; Value: string }[],
  labels: string[]
}

export interface rootRouteObject {
  URL: string,
  Body: { Name: string; Value: string }[],
  Params: { Name: string; Value: string }[],
  Header: { OfType: string; Value: string }[],
  Method: string,
  Labels: string[]
}

export interface paramsTransformValue {
  key: string;
  value: string
}

export interface paramsObject {
  Name: string;
  Value: string;
}

export interface bodyObject {
  Name: string;
  Value: string;
}

export interface headersObject {
  OfType: string;
  Value: string
}

export interface LabelType {
  name: string;
  color: string;
}
