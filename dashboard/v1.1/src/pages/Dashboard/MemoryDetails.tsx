import React, { FC } from 'react';
import Chart from 'react-apexcharts';
import { chartData } from '../../utils/queryTypes';

interface MemoryDetailsProps {
  availableBytes: chartData[];
  freeBytes: chartData[];
  totalBytes: chartData[];
  usedBytes: chartData[];
  darkMode(status: boolean): any;
}

const MemoryDetails: FC<MemoryDetailsProps> = ({
  availableBytes,
  freeBytes,
  totalBytes,
  usedBytes,
  darkMode
}) => {
  const series = [
    {
      name: 'Available',
      data: availableBytes
    },
    {
      name: 'Free',
      data: freeBytes
    },
    {
      name: 'Total',
      data: totalBytes
    },
    {
      name: 'Used',
      data: usedBytes
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
      text: 'Memory (RAM) details (in kilo-bytes)',
      align: 'center'
    },
    fill: {
      opacity: 1,
      type: 'gradient',
      gradient: {
        shade: 'dark',
        type: 'vertical',
        shadeIntensity: 0.1,
        inverseColors: true,
        opacityFrom: 0.8,
        opacityTo: 0.2
      }
    },
    tooltip: {
      theme: !darkMode ? 'light' : 'dark'
    }
  };

  return (
    <>
      <Chart series={series} options={options} height="500" type="area" />
    </>
  );
};

export default React.memo(MemoryDetails);
