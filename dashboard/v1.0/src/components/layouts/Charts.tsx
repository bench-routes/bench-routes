import React, { FC } from 'react';
import { Line } from 'react-chartjs-2';
import { Datasets } from '../../utils/types';

interface ChartProps {
  opts: ChartOptions[];
}

interface ChartFormattedOptions {
  datasets: Datasets[];
  labels: any[];
}

export interface ChartOptions {
  xAxisValues: any[];
  yAxisValues: any[];
  fill: boolean;
  lineTension: number;
  backgroundColor: string;
  borderColor: string;
  pBorderColor: string;
  pBackgroundColor: string;
  pBorderWidth: number;
  pHoverBorderWidth: number;
  pHoverRadius: number;
  pRadius: number;
  label: string;
  pHoverBackgroundColor: string;
  pHoverBorderColor: string;
}

export const ChartValues = (
  xAxisValues?: any,
  yAxisValues?: any,
  label?: any,
  colorCode?: any
): ChartOptions => {
  // If default values are required, then it can be called as `ChartValues();`
  if (!xAxisValues && !yAxisValues && !label && !colorCode) {
    return {
      backgroundColor: '',
      borderColor: '',
      fill: false,
      label: '',
      lineTension: 0.1,
      pBackgroundColor: '#fff',
      pBorderColor: 'rgba(75,192,2,1)',
      pBorderWidth: 1,
      pHoverBackgroundColor: 'rgba(7,12,19,0.4)',
      pHoverBorderColor: 'rgba(220,220,220,1)',
      pHoverBorderWidth: 2,
      pHoverRadius: 5,
      pRadius: 1,
      xAxisValues: [],
      yAxisValues: []
    };
  }
  return {
    backgroundColor: colorCode,
    borderColor: colorCode,
    fill: false,
    label,
    lineTension: 0.1,
    pBackgroundColor: '#fff',
    pBorderColor: 'rgba(75,192,2,1)',
    pBorderWidth: 1,
    pHoverBackgroundColor: 'rgba(7,12,19,0.4)',
    pHoverBorderColor: 'rgba(220,220,220,1)',
    pHoverBorderWidth: 2,
    pHoverRadius: 5,
    pRadius: 1,
    xAxisValues,
    yAxisValues
  };
};

export const Charts: FC<ChartProps> = ({ opts }) => {
  const formatProps = (options: ChartOptions[]): ChartFormattedOptions => {
    const data: Datasets[] = [];
    for (const i of options) {
      data.push({
        backgroundColor: i.backgroundColor,
        borderColor: i.borderColor,
        data: i.yAxisValues,
        fill: i.fill,
        label: i.label,
        lineTension: i.lineTension,
        pointBackgroundColor: i.pBackgroundColor,
        pointBorderColor: i.pBorderColor,
        pointBorderWidth: i.pBorderWidth,
        pointHoverBackgroundColor: i.pBackgroundColor,
        pointHoverBorderColor: i.pHoverBorderColor,
        pointHoverBorderWidth: i.pHoverBorderWidth,
        pointHoverRadius: i.pHoverRadius,
        pointRadius: i.pRadius
      });
    }

    return {
      datasets: data,
      labels: options[0].xAxisValues
    };
  };

  const options = formatProps(opts);
  const wrapper = 'canvas-chart-wrapper';

  return (
    <div className={wrapper} style={{ height: '100%' }}>
      {options ? <Line data={options} /> : null}
    </div>
  );
};
