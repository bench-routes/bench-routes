import React, { Component } from 'react';
import BRConnect from '../../utils/connection';
import { ChartOptions, Charts, ChartValues } from '../layouts/Charts';
import Submenu from '../layouts/Submenu';

interface FPingModulePropsTypes {}

interface FPingModuleStateTypes {
  routes: object;
  sAddress: string;
  chartOpts: ChartOptions[];
  showChart: boolean;
  packetLossChartOpts: ChartOptions[];
}

export default class FloodPing extends Component<
  FPingModulePropsTypes,
  FPingModuleStateTypes
> {
  private connection: BRConnect;
  constructor(props: FPingModulePropsTypes) {
    super(props);

    this.connection = new BRConnect();
    const tmp: ChartOptions[] = [ChartValues()];
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
        const norTime: number[] = [];
        const timeStamp: string[] = [];
        const yMin: ChartOptions[] = [];
        const yMean: ChartOptions[] = [];
        const yMax: ChartOptions[] = [];
        const yMdev: ChartOptions[] = [];
        const packetLoss: number[] = [];

        let inst;
        for (inst of data) {
          yMin.push(inst.Min);
          yMean.push(inst.Mean);
          yMax.push(inst.Max);
          yMdev.push(inst.Mdev);
          norTime.push(inst.relative);
          timeStamp.push(inst.timestamp);
          packetLoss.push(inst.PacketLoss);
        }

        const options: ChartOptions[] = [
          ChartValues(norTime, yMin, 'Minimum', 'rgba(75,192,192,0.4)'),
          ChartValues(norTime, yMean, 'Mean', 'rgba(75,192,2,0.4)'),
          ChartValues(norTime, yMax, 'Maximum', 'rgba(5,192,19,0.4)'),
          ChartValues(norTime, yMdev, 'Standard-Deviation', 'rgba(7,12,19,0.4)')
        ];

        const optionsPacketLoss: ChartOptions[] = [
          ChartValues(
            norTime,
            packetLoss,
            'Packet-loss',
            'rgba(75,192,192,0.4)'
          )
        ];

        this.setState({
          chartOpts: options,
          packetLossChartOpts: optionsPacketLoss,
          showChart: true
        });
      });
  };

  public opts = (operation: string) => {
    switch (operation) {
      case 'start':
        this.connection.signalFloodPingStart().then(res => {
          if (res.data) {
            alert('Flood Ping routine started');
          }
        });
        break;
      case 'stop':
        this.connection.signalFloodPingStop().then(res => {
          if (res.data) {
            alert('Flood Ping routine stopped');
          }
        });
        break;
    }
  };

  public render() {
    return (
      <>
        <div className="btn-layout">
          {/* operations */}
          <button
            className="button-operations btn btn-success"
            onClick={() => this.opts('start')}
          >
            Start
          </button>
          <button
            className="button-operations btn btn-danger"
            onClick={() => this.opts('start')}
          >
            Stop
          </button>
        </div>
        <Submenu
          module="ping"
          submodule=""
          getAddress={this.getAddressSubmenu}
        />
        {this.state.showChart ? (
          <div style={{ overflow: 'scroll', height: '45%' }}>
            <Charts opts={this.state.chartOpts} />
            <Charts opts={this.state.packetLossChartOpts} />
          </div>
        ) : (
          <div>Chart not available</div>
        )}
      </>
    );
  }
}
