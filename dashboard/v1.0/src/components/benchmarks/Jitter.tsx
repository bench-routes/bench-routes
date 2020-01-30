import React, { Component } from 'react';
import BRConnect from '../../utils/connection';
import { ChartOptions, Charts, ChartValues } from '../layouts/Charts';
import Submenu from '../layouts/Submenu';

interface JitterModulePropsTypes {}

interface JitterModuleStateTypes {
  routes: object;
  sAddress: string;
  chartOpts: ChartOptions[];
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
    const tmp: ChartOptions[] = [ChartValues()];
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
      .signalJitterRouteFetchAllTimeSeries(sAddressParam)
      .then((res: any) => {
        const data: any[] = JSON.parse(res.data);
        const jitter: ChartOptions[] = [];
        const norTime: number[] = [];
        const timeStamp: string[] = [];

        let inst;
        for (inst of data) {
          jitter.push(inst.datapoint);
          norTime.push(inst.relative);
          timeStamp.push(inst.timestamp);
        }

        const options: ChartOptions[] = [
          ChartValues(norTime, jitter, 'Jitter', 'rgba(75,192,192,0.4)')
        ];

        this.setState({ chartOpts: options, showChart: true });
      });
  };

  public opts = (operation: string) => {
    switch (operation) {
      case 'start':
        this.connection.signalJitterStart().then(res => {
          if (res.data) {
            alert('Jitter routine started');
          }
        });
        break;
      case 'stop':
        this.connection.signalJitterStop().then(res => {
          if (res.data) {
            alert('Jitter routine stopped');
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
            onClick={() => this.opts('stop')}
          >
            Stop
          </button>
        </div>
        <Submenu
          module="jitter"
          submodule=""
          getAddress={this.getAddressSubmenu}
        />
        {this.state.showChart ? (
          <div>
            <Charts opts={this.state.chartOpts} />
          </div>
        ) : (
          <div>Chart not available</div>
        )}
      </>
    );
  }
}
