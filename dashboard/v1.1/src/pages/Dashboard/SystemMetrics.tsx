import React, { FC } from 'react';
import { useFetch } from '../../utils/useFetch';
import CPUUsage from './CPUUsage';
import MemoryUsagePercent from './MemoryUsage';
import DiskUsage from './Disk';
import TimeInstance from '../../utils/brt';

export interface queryValueCPUUsage {
  CPUUsage: string;
  normalizedTime: number;
}

export interface queryValueMemoryUsedPercent {
  memoryUsedPercent: string;
  normalizedTime: number;
}

interface QueryRange {
  start: number;
  end: number;
}

export interface queryValueDisk {
  cached: string;
  diskIO: string;
  normalizedTime: number;
}

interface queryValueMemory {
  availableBytes: string;
  freeBytes: string;
  totalBytes: string;
  usedBytes: string;
  usedPercent: string;
  normalizedTime: number;
}

interface queryValueSystemMetrics {
  cpuTotalUsage: string;
  disk: queryValueDisk;
  memory: queryValueMemory;
}

interface QueryValues {
  normalizedTime: number;
  timestamp: string;
  value: queryValueSystemMetrics;
}

interface QueryResponse {
  queryTime: number;
  range: QueryRange;
  timeSeriesPath: string;
  values: QueryValues[];
}

const url = 'http://localhost:9090';

const segregateMetrics = (metricValues: QueryValues[]) => {
  const cpuUsageSlice: queryValueCPUUsage[] = [];
  const diskSlice: queryValueDisk[] = [];
  const memorySlice: queryValueMemory[] = [];
  const memoryUsedPercentSlice: queryValueMemoryUsedPercent[] = [];

  let prev = metricValues[metricValues.length - 1].normalizedTime;
  for (const metric of metricValues) {
    cpuUsageSlice.push({
      CPUUsage: metric.value.cpuTotalUsage,
      normalizedTime: prev - metric.normalizedTime
    });
    diskSlice.push({
      cached: metric.value.disk.cached,
      diskIO: metric.value.disk.diskIO,
      normalizedTime: prev - metric.normalizedTime
    });
    memorySlice.push({
      availableBytes: metric.value.memory.availableBytes,
      freeBytes: metric.value.memory.freeBytes,
      totalBytes: metric.value.memory.totalBytes,
      usedBytes: metric.value.memory.usedBytes,
      usedPercent: metric.value.memory.usedPercent,
      normalizedTime: prev - metric.normalizedTime
    });
    memoryUsedPercentSlice.push({
      memoryUsedPercent: metric.value.memory.usedPercent,
      normalizedTime: prev - metric.normalizedTime
    });
  }
  return { cpuUsageSlice, diskSlice, memorySlice, memoryUsedPercentSlice };
};

const SystemMetrics: FC<{}> = () => {
  const endTimestamp = new Date().getTime() * 1000000 - TimeInstance.Hour;
  const { response, error, isLoading } = useFetch<QueryResponse>(`${url}/query?timeSeriesPath=storage/system&endTimestamp=${endTimestamp}`);
  if (error) {
    console.warn(error)
  }
  if (!response.data) {
    return <>loading...</>
  }

  const responseInFormat = segregateMetrics(response.data.values);

  return (
    <div className="column">
      <div className="row">
        <div className="col-md-6">
          <CPUUsage cpuMetrics={responseInFormat.cpuUsageSlice} />
        </div>
        <div className="col-md-6">
          <MemoryUsagePercent memoryUsagePercentMetrics={responseInFormat.memoryUsedPercentSlice} />
        </div>
        <div className="col-md-12">
          <DiskUsage metrics={responseInFormat.diskSlice} />
        </div>
      </div>
    </div>
  );
};

export default SystemMetrics;