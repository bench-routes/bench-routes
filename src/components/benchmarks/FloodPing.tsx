import React, { Component } from 'react';
import BRConnect from '../../utils/connection';
import {
  BRChartOpts,
  BRCharts,
  chartDefaultOptValues
} from '../layouts/BRCharts';
import Submenu from '../layouts/Submenu';

interface FPingModulePropsTypes {}

interface FPingModuleStateTypes {
  routes: object;
  sAddress: string;
  chartOpts: BRChartOpts[];
  showChart: boolean;
  packetLossChartOpts: BRChartOpts[];
}

export default class FloodPing extends Component<
  FPingModulePropsTypes,
  FPingModuleStateTypes
> {
  private connection: BRConnect;
  constructor(props: FPingModulePropsTypes) {
    super(props);

    this.connection = new BRConnect();
    const tmp: BRChartOpts[] = [chartDefaultOptValues];
    this.state = {
      // submenu address
      chartOpts: tmp,
      packetLossChartOpts: tmp,
      routes: {},
      sAddress: '',
      showChart: false
    };
  }

  public getAddressSubmenu = (sAddressParam: string) => {
    this.setState({ sAddress: sAddressParam, showChart: false });
    this.connection
      .signalFloodPingRouteFetchAllTimeSeries(sAddressParam)
      .then((res: any) => {
        const data: any[] = JSON.parse(res.data);
        const yMin: BRChartOpts[] = [];
        const yMean: BRChartOpts[] = [];
        const yMax: BRChartOpts[] = [];
        const yMdev: BRChartOpts[] = [];
        const norTime: number[] = [];
        const timeStamp: string[] = [];
        const packetLoss: number[] = [];

        let inst;
        for (inst of data) {
          yMin.push(inst.datapoint.Min);
          yMean.push(inst.datapoint.Mean);
          yMax.push(inst.datapoint.Max);
          yMdev.push(inst.datapoint.Mdev);
          norTime.push(inst.normalizedTime);
          timeStamp.push(inst.timestamp);
          packetLoss.push(inst.datapoint.PacketLoss);
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

        const packetLossChartOptions: BRChartOpts[] = [
          {
            backgroundColor: 'rgba(75,192,192,0.4)',
            borderColor: 'rgba(255,20,147,1)',
            fill: false,
            label: 'Packet-Loss',
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
            yAxisValues: packetLoss
          }
        ];

        this.setState({
          chartOpts: chartOptions,
          packetLossChartOpts: packetLossChartOptions,
          showChart: true
        });
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
          <div style={{ overflow: 'auto' }}>
            <div>
              <BRCharts opts={this.state.chartOpts} />
            </div>
            <div>
              <BRCharts opts={this.state.packetLossChartOpts} />
            </div>
          </div>
        ) : (
          <div>Chart not available</div>
        )}
      </>
    );
  }
}
