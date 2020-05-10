import React, { FC } from 'react';
import Chart from 'react-apexcharts';
import { chartData } from '../../utils/queryTypes';

interface MemoryDetailsProps {
  availableBytes: chartData[];
  freeBytes: chartData[];
  totalBytes: chartData[];
  usedBytes: chartData[];
}

const MemoryDetails: FC<MemoryDetailsProps> = ({
  availableBytes,
  freeBytes,
  totalBytes,
  usedBytes
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

export default React.memo(MemoryDetails);
