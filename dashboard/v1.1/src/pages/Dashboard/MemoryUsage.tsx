import React, { FC } from 'react';
import Chart from 'react-apexcharts';
import { chartData } from '../../utils/queryTypes';

interface MemoryUsagePercentProps {
  memoryUsagePercentMetrics: chartData[];
}

const MemoryUsagePercent: FC<MemoryUsagePercentProps> = ({
  memoryUsagePercentMetrics
}) => {
  const dataFormatted = memoryUsagePercentMetrics;
  const series = [
    {
      name: 'Memory',
      data: dataFormatted
    }
  ];
  const options = {
    chart: {
      type: 'area',
      background: '#fff'
    },
    dataLabels: {
      enabled: false
    },
    stroke: {
      show: true,
      curve: 'straight',
      lineCap: 'butt',
      width: 1
    },
    subtitle: {
      text: 'Memory (RAM) usage  (in percent)',
      align: 'center'
    },
    fill: {
      opacity: 1,
      type: 'gradient',
      gradient: {
        shade: 'dark',
        type: 'vertical',
        shadeIntensity: 0.6,
        inverseColors: true,
        opacityFrom: 1,
        opacityTo: 0.2
      }
    }
  };

  return (
    <>
      <Chart series={series} options={options} height="300" type="area" />
    </>
  );
};

export default React.memo(MemoryUsagePercent);
