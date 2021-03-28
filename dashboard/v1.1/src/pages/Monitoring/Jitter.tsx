import React, { FC, useContext } from 'react';
import Chart from 'react-apexcharts';
import Alert from '@material-ui/lab/Alert';
import { ThemeContext, XticksContext } from '../../layouts/BaseLayout';
import { chartData } from '../../utils/queryTypes';

interface JitterProps {
  value: chartData[];
}

const Jitter: FC<JitterProps> = ({ value }) => {
  const themeMode = useContext(ThemeContext);
  const xticks = useContext(XticksContext);
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
      }
    },
    yaxis: {
      title: {
        text: 'milliseconds'
      }
    },
    xaxis: {
      title: {
        text: 'Time'
      },
      tickAmount: Number(xticks)
    },
    theme: {
      mode: themeMode
    }
  };
  return (
    <>
      {!value.length ? (
        <Alert severity="error">No data found</Alert>
      ) : (
        <Chart series={series} options={options} height="300" />
      )}
    </>
  );
};

export default Jitter;
