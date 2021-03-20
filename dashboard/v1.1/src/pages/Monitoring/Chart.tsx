import React from 'react';
import Chart, { Props } from 'react-apexcharts';
import Alert from '@material-ui/lab/Alert';

import { chartData } from '../../utils/queryTypes';
import { getChartSeries } from '../../utils/chart';

interface ChartProps extends Props {
  values: chartData[];
  name: string;
}

const ChartComponent = ({ options, values, name }: ChartProps) => {
  const { isValid, series } = getChartSeries(name, values);

  return (
    <>
      {isValid ? (
        <Alert severity="error">No data found</Alert>
      ) : (
        <Chart series={series} options={options} height="300" />
      )}
    </>
  );
};

export default ChartComponent;
