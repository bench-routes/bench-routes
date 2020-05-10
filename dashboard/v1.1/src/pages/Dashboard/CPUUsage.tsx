import React, { FC } from 'react';
import Chart from 'react-apexcharts';
import { chartData } from '../../utils/queryTypes';

interface CPUUsageProps {
  cpuMetrics: chartData[];
}

const CPUUsage: FC<CPUUsageProps> = ({ cpuMetrics }) => {
  const dataFormatted = cpuMetrics;
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
      text: 'CPU usage (in percent)',
      align: 'center'
    }
  };

  return (
    <>
      <Chart series={series} options={options} height="300" />
    </>
  );
};

export default React.memo(CPUUsage);
