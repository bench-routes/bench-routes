import React, { FC } from 'react';
import Chart from 'react-apexcharts';
import { chartData } from '../../utils/queryTypes';

interface JitterProps {
  value: chartData[];
}

const Jitter: FC<JitterProps> = ({ value }) => {
  const series = [
    {
      name: 'min',
      data: value
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
    }
  };
  return <Chart series={series} options={options} height="300" />;
};

export default Jitter;
