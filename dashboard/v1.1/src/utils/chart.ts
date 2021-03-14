import { chartData } from './queryTypes';

export const getChartSeries = (name: string, values: chartData[]) => {
  let series = [{}];
  let isValid = false;

  switch (name) {
    case 'Ping':
      series = [
        {
          name: 'min',
          data: [values[0]]
        },
        {
          name: 'mean',
          data: [values[1]]
        },
        {
          name: 'max',
          data: [values[2]]
        }
      ];
      if (!values.length) isValid = true;
      return { series, isValid };
    default:
      series = [
        {
          name: name,
          data: values
        }
      ];
      if (!values.length) isValid = true;
      return { series, isValid };
  }
};
