import React, { FC } from 'react';
import Chart from 'react-apexcharts';
import { queryValueMemoryUsedPercent } from '../../utils/queryTypes';

interface chartData {
  x: number;
  y: string;
}

const chartFormating = (metrics: queryValueMemoryUsedPercent[]) => {
  const chartData: chartData[] = [];
  for (const metric of metrics) {
    chartData.push({
      y: metric.memoryUsedPercent,
      x: metric.normalizedTime
    });
  }
  return chartData;
};

interface MemoryUsagePercentProps {
  memoryUsagePercentMetrics: queryValueMemoryUsedPercent[];
}

const MemoryUsagePercent: FC<MemoryUsagePercentProps> = ({
  memoryUsagePercentMetrics
}) => {
  const dataFormatted = chartFormating(memoryUsagePercentMetrics);
  const series = [
    {
      name: 'Memory',
      data: dataFormatted
    }
  ];
  const options = {
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
      text: 'Memory (RAM) usage  (in percent)',
      align: 'center'
    }
  };

  return (
    <>
      <Chart series={series} options={options} height="300" />
    </>
  );
};

export default MemoryUsagePercent;
