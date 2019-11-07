import React from 'react';
import BRConnect from '../../utils/connection';
import {
  BRChartOpts,
  BRCharts,
  chartDefaultOptValues
} from '../layouts/BRCharts';
import Submenu from '../layouts/Submenu';

interface PingModulePropsTypes {}

interface PingModuleStateTypes {
  routes: object;
  sAddress: string;
  chartOpts: BRChartOpts[];
  showChart: boolean;
}

export default class PingModule extends React.Component<
  PingModulePropsTypes,
  PingModuleStateTypes
> {
  private connection: BRConnect;
  constructor(props: PingModulePropsTypes) {
    super(props);

    this.connection = new BRConnect();
    const tmp: BRChartOpts[] = [chartDefaultOptValues];
    this.state = {
      chartOpts: tmp,
      routes: {},
      sAddress: '',
      showChart: false
    };
  }

  public getAddressSubmenu = (sAddressParam: string) => {
    this.setState({ sAddress: sAddressParam, showChart: false });
    this.connection
      .signalPingRouteFetchAllTimeSeries(sAddressParam)
      .then((res: any) => {
        // format the values
        const data: any[] = JSON.parse(res.data);
        const yMin: BRChartOpts[] = [];
        const yMean: BRChartOpts[] = [];
        const yMax: BRChartOpts[] = [];
        const yMdev: BRChartOpts[] = [];
        const norTime: number[] = [];
        const timeStamp: string[] = [];

        // tslint:disable-next-line: prefer-for-of
        for (let i = 0; i < data.length; i++) {
          const inst: any = data[i];
          yMin.push(inst.datapoint.Min);
          yMean.push(inst.datapoint.Mean);
          yMax.push(inst.datapoint.Max);
          yMdev.push(inst.datapoint.Mdev);
          norTime.push(inst.normalizedTime);
          timeStamp.push(inst.timestamp);
        }

        const chartOptions: BRChartOpts[] = [
          {
            backgroundColor: 'rgba(75,192,192,0.4)',
            borderColor: 'rgba(75,192,192,1)',
            fill: false,
            label: 'Minimum',
            lineTension: 0.1,
            pBackgroundColor: '#fff',
            pBorderColor: 'rgba(75,192,2,1)',
            pBorderWidth: 1,
            pHoverBackgroundColor: 'rgba(75,192,192,1)',
            pHoverBorderColor: 'rgba(220,220,220,1)',
            pHoverBorderWidth: 2,
            pHoverRadius: 5,
            pRadius: 1,
            xAxisValues: norTime,
            yAxisValues: yMin
          },
          {
            backgroundColor: 'rgba(75,192,2,0.4)',
            borderColor: 'rgba(75,192,2,0.4)',
            fill: false,
            label: 'Mean',
            lineTension: 0.1,
            pBackgroundColor: '#fff',
            pBorderColor: 'rgba(75,192,2,1)',
            pBorderWidth: 1,
            pHoverBackgroundColor: 'rgba(75,2,192,1)',
            pHoverBorderColor: 'rgba(220,220,220,1)',
            pHoverBorderWidth: 2,
            pHoverRadius: 5,
            pRadius: 1,
            xAxisValues: norTime,
            yAxisValues: yMean
          },
          {
            backgroundColor: 'rgba(5,192,19,0.4)',
            borderColor: 'rgba(5,192,19,0.4)',
            fill: false,
            label: 'Maximum',
            lineTension: 0.1,
            pBackgroundColor: '#fff',
            pBorderColor: 'rgba(75,192,2,1)',
            pBorderWidth: 1,
            pHoverBackgroundColor: 'rgba(75,192,192,1)',
            pHoverBorderColor: 'rgba(220,220,220,1)',
            pHoverBorderWidth: 2,
            pHoverRadius: 5,
            pRadius: 1,
            xAxisValues: norTime,
            yAxisValues: yMax
          },
          {
            backgroundColor: 'rgba(7,12,19,0.4)',
            borderColor: 'rgba(7,12,19,0.4)',
            fill: false,
            label: 'Minimum-deviation',
            lineTension: 0.1,
            pBackgroundColor: '#fff',
            pBorderColor: 'rgba(75,192,2,1)',
            pBorderWidth: 1,
            pHoverBackgroundColor: 'rgba(7,12,19,0.4)',
            pHoverBorderColor: 'rgba(220,220,220,1)',
            pHoverBorderWidth: 2,
            pHoverRadius: 5,
            pRadius: 1,
            xAxisValues: norTime,
            yAxisValues: yMdev
          }
        ];

        this.setState({ chartOpts: chartOptions, showChart: true });
      });
  };

  public render() {
    return (
      <>
        <Submenu
          module="ping"
          submodule=""
          getAddress={this.getAddressSubmenu}
        />
        {this.state.showChart ? (
          <div>
            <BRCharts opts={this.state.chartOpts} />
          </div>
        ) : (
          <div>Chart not available</div>
        )}
      </>
    );
  }
}
