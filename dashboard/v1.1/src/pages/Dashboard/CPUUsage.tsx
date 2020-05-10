import React, { FC } from 'react';
import Chart from 'react-apexcharts';
import { queryValueCPUUsage } from '../../utils/queryTypes';

interface chartData {
  x: number;
  y: string;
}

const chartFormating = (metrics: queryValueCPUUsage[]) => {
  const chartData: chartData[] = [];

  for (const metric of metrics) {
    chartData.push({
      y: metric.CPUUsage,
      x: metric.normalizedTime
    });
  }

  return chartData;
};

interface CPUUsageProps {
  cpuMetrics: queryValueCPUUsage[];
}

const CPUUsage: FC<CPUUsageProps> = ({ cpuMetrics }) => {
  const dataFormatted = chartFormating(cpuMetrics);
  const series = [
    {
      name: 'CPU',
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
      text: 'CPU usage',
      align: 'center'
    }
  };

  return (
    <>
      <Chart series={series} options={options} height="300" />
    </>
  );
};

export default CPUUsage;
