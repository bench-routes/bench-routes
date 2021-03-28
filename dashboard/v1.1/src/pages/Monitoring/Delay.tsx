import React, { FC, useContext } from 'react';
import Chart from 'react-apexcharts';
import Alert from '@material-ui/lab/Alert';
import { ThemeContext } from '../../layouts/BaseLayout';
import { chartData } from '../../utils/queryTypes';

interface DelayProps {
  delay: chartData[];
}

const Delay: FC<DelayProps> = ({ delay }) => {
  const themeMode = useContext(ThemeContext);
  const series = [
    {
      name: 'Delay',
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
      }
    },
    yaxis: {
      title: {
        text: 'Response (in ms)'
      }
    },
    xaxis: {
      title: {
        text: 'Time'
      }
    },
    theme: {
      mode: themeMode
    }
  };
  return (
    <>
      {!delay.length ? (
        <Alert severity="error">No data found</Alert>
      ) : (
        <Chart series={series} options={options} height="300" />
      )}
    </>
  );
};

export default Delay;
