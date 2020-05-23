import React, { FC } from 'react';
import Chart from 'react-apexcharts';
import { chartData } from '../../utils/queryTypes';

interface DelayProps {
  delay: chartData[];
}

const Delay: FC<DelayProps> = ({ delay }) => {
  const series = [
    {
      name: 'CPU',
      data: delay
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

export default Delay;
