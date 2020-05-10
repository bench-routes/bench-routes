import React, { FC } from 'react';
import Chart from 'react-apexcharts';
import { queryValueDisk } from '../../utils/queryTypes';

interface chartData {
  x: number;
  y: string;
}

const chartFormating = (metrics: queryValueDisk[]) => {
  const chartDataDiskIO: chartData[] = [];
  const chartDataCache: chartData[] = [];

  for (const metric of metrics) {
    chartDataDiskIO.push({
      y: metric.diskIO,
      x: metric.normalizedTime
    });
    chartDataCache.push({
      y: metric.cached,
      x: metric.normalizedTime
    });
  }

  return { chartDataDiskIO, chartDataCache };
};

interface DiskUsageProps {
  metrics: queryValueDisk[];
}

const DiskUsage: FC<DiskUsageProps> = ({ metrics }) => {
  const { chartDataDiskIO, chartDataCache } = chartFormating(metrics);

  const seriesDiskIO = [
    {
      name: 'Disk IO in bytes (+ve means write / -ve means read)',
      data: chartDataDiskIO
    }
  ];
  const seriesCache = [
    {
      name: 'Cache (in bytes)',
      data: chartDataCache
    }
  ];

  const optionsDiskIO = {
    chart: {
      type: 'area',
      animations: {
        enabled: true,
        easing: 'easeinout',
        speed: 800,
        animateGradually: {
          enabled: true,
          delay: 150
        },
        dynamicAnimation: {
          enabled: true,
          speed: 350
        }
      },
      background: '#fff'
    },
    datalabels: {
      enabled: false
    },
    stroke: {
      show: true,
      curve: 'straight',
      lineCap: 'butt',
      width: 3
    },
    subtitle: {
      text: 'Disk IO in bytes (+ve means write / -ve means read)',
      align: 'center'
    }
  };
  const optionsCache = {
    chart: {
      type: 'area',
      animations: {
        enabled: true,
        easing: 'easeinout',
        speed: 800,
        animateGradually: {
          enabled: true,
          delay: 150
        },
        dynamicAnimation: {
          enabled: true,
          speed: 350
        }
      },
      background: '#fff'
    },
    datalabels: {
      enabled: false
    },
    stroke: {
      show: true,
      curve: 'straight',
      lineCap: 'butt',
      width: 3
    },
    subtitle: {
      text: 'Cache (in bytes)',
      align: 'center'
    }
  };

  return (
    <div className="row">
      <div className="col-md-6">
        <Chart series={seriesDiskIO} options={optionsDiskIO} height="300" />
      </div>
      <div className="col-md-6">
        <Chart series={seriesCache} options={optionsCache} height="300" />
      </div>
    </div>
  );
};

export default DiskUsage;
