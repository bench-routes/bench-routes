import React from 'react';
import { Line } from 'react-chartjs-2';

export interface BRChartOpts {
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

// override xAxisValues, yAxisValues, label, backgroundColor, borderColor, pHoverBackgroundColor, pHoverBorderColor
export let chartDefaultOptValues: BRChartOpts = {
  xAxisValues: [],
  yAxisValues: [],
  fill: false,
  lineTension: 0.1,
  backgroundColor: 'rgba(75,192,192,0.4)',
  borderColor: 'rgba(75,192,192,1)',
  pBackgroundColor: '#fff',
  pBorderColor: 'rgba(75,192,2,1)',
  pBorderWidth: 1,
  pHoverBorderWidth: 2,
  pHoverRadius: 5,
  pRadius: 1,
  label: '',
  pHoverBackgroundColor: 'rgba(75,192,192,1)',
  pHoverBorderColor: 'rgba(220,220,220,1)'
};

interface BRChartProps {
  opts: BRChartOpts[];
}

interface Datasets {
  label: string;
  fill: boolean;
  lineTension: number;
  backgroundColor: string;
  borderColor: string;
  pointBorderColor: string;
  pointBackgroundColor: string;
  pointBorderWidth: number;
  pointHoverRadius: number;
  pointHoverBackgroundColor: string;
  pointHoverBorderColor: string;
  pointHoverBorderWidth: number;
  pointRadius: number;
  data: any[];
}

interface BRChartOptionsFormated {
  labels: any[];
  datasets: Datasets[];
}

interface BRChartState {
  wrapperClassName: string;
  chartJSOpts: any;
}

export class BRCharts extends React.Component<BRChartProps, BRChartState> {
  constructor(props: BRChartProps) {
    super(props);
    this.state = {
      chartJSOpts: null,
      wrapperClassName: 'canvas-chart-wrapper'
    };
  }

  public componentDidMount() {
    this.setState({ chartJSOpts: this.formatProps(this.props.opts) });
  }

  public render() {
    return (
      <div className={this.state.wrapperClassName} style={{ height: '100%' }}>
        {this.state.chartJSOpts ? <Line data={this.state.chartJSOpts} /> : null}
      </div>
    );
  }

  private formatProps = (inst: BRChartOpts[]) => {
    const data: Datasets[] = [];
    for (const i of inst) {
      const tmp: Datasets = {
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
      };
      data.push(tmp);
    }
    return {
      datasets: data,
      labels: inst[0].xAxisValues
    };
  };
}
