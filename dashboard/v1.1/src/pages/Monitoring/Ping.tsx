import React, { FC } from 'react';
import Chart from 'react-apexcharts';
import { chartData } from '../../utils/queryTypes';

interface PingProps {
  min: chartData[];
  mean: chartData[];
  max: chartData[];
}

const Ping: FC<PingProps> = ({ min, mean, max }) => {
  const series = [
    {
      name: 'min',
      data: min
    },
    {
      name: 'mean',
      data: mean
    },
    {
      name: 'max',
      data: max
    }
  ];
  const options = {
    chart: {
      type: 'area',
      xaxis: {
        type: 'category',
        categories: [],
        labels: {
          show: true,
          rotate: 0,
          rotateAlways: false,
          hideOverlappingLabels: true,
          trim: true
        }
      },
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

export default Ping;
