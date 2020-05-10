import React, { FC } from 'react';
import Chart from 'react-apexcharts';
import { queryValueMemory } from '../../utils/queryTypes';

interface chartData {
  x: number;
  y: string;
}

const chartFormating = (metrics: queryValueMemory[]) => {
  const chartDataAvailableBytes: chartData[] = [];
  const chartDataFreeBytes: chartData[] = [];
  const chartDataTotalBytes: chartData[] = [];
  const chartDataUsedBytes: chartData[] = [];

  for (const metric of metrics) {
    chartDataAvailableBytes.push({
      y: metric.availableBytes,
      x: metric.normalizedTime
    });
    chartDataFreeBytes.push({
      y: metric.freeBytes,
      x: metric.normalizedTime
    });
    chartDataTotalBytes.push({
      y: metric.totalBytes,
      x: metric.normalizedTime
    });
    chartDataUsedBytes.push({
      y: metric.usedBytes,
      x: metric.normalizedTime
    });
  }

  return {
    chartDataAvailableBytes,
    chartDataFreeBytes,
    chartDataTotalBytes,
    chartDataUsedBytes
  };
};

interface MemoryDetailsProps {
  metrics: queryValueMemory[];
}

const MemoryDetails: FC<MemoryDetailsProps> = ({ metrics }) => {
  const dataFormatted = chartFormating(metrics);
  const series = [
    {
      name: 'Available',
      data: dataFormatted.chartDataAvailableBytes
    },
    {
      name: 'Free',
      data: dataFormatted.chartDataFreeBytes
    },
    {
      name: 'Total',
      data: dataFormatted.chartDataTotalBytes
    },
    {
      name: 'Used',
      data: dataFormatted.chartDataUsedBytes
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
      text: 'Memory (RAM) details (in kilo-bytes)',
      align: 'center'
    }
  };

  return (
    <>
      <Chart series={series} options={options} height="500" />
    </>
  );
};

export default MemoryDetails;
