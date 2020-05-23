import React, { FC } from 'react';
import Chart from 'react-apexcharts';
import { chartData } from '../../utils/queryTypes';

interface ResLengthProps {
  resLength: chartData[];
}

const ResLength: FC<ResLengthProps> = ({ resLength }) => {
  const series = [
    {
      name: 'CPU',
      data: resLength
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

export default ResLength;
