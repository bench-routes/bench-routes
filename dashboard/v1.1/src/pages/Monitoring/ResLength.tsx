import React, { FC } from 'react';
import Chart from 'react-apexcharts';
import Alert from '@material-ui/lab/Alert';

import { chartData } from '../../utils/queryTypes';

interface ResLengthProps {
  resLength: chartData[];
}

const ResLength: FC<ResLengthProps> = ({ resLength }) => {
  const series = [
    {
      name: 'Response length',
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
    },
    yaxis: {
      title: {
        text: 'Length'
      }
    },
    xaxis: {
      title: {
        text: 'Time'
      }
    }
  };
  return (
    <>
      {!resLength.length ? (
        <Alert severity="error">No data found</Alert>
      ) : (
        <Chart series={series} options={options} height="300" />
      )}
    </>
  );
};

export default ResLength;
