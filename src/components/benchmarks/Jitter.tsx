import React, { Component } from 'react';
import BRConnect from '../../utils/connection';
import {
  BRChartOpts,
  BRCharts,
  chartDefaultOptValues
} from '../layouts/BRCharts';
import Submenu from '../layouts/Submenu';

interface JitterModulePropsTypes {}

interface JitterModuleStateTypes {
  routes: object;
  sAddress: string;
  chartOpts: BRChartOpts[];
  showChart: boolean;
}

export default class Jitter extends Component<
  JitterModulePropsTypes,
  JitterModuleStateTypes
> {
  private connection: BRConnect;
  constructor(props: JitterModulePropsTypes) {
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
            label: 'Jitter',
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
          }
        ];

        this.setState({ chartOpts: chartOptions, showChart: true });
      });
  };

  // public getAddressSubmenu = (sAddressParam: string) => {
  //   this.setState({ sAddress: sAddressParam });
  // };

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
