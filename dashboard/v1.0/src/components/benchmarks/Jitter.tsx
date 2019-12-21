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
      .signalJitterRouteFetchAllTimeSeries(sAddressParam)
      .then((res: any) => {
        // format the values
        const data: any[] = JSON.parse(res.data);
        const jitter: BRChartOpts[] = [];
        const norTime: number[] = [];
        const timeStamp: string[] = [];

        let inst;
        for (inst of data) {
          jitter.push(inst.datapoint);
          norTime.push(inst.relative);
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
            yAxisValues: jitter
          }
        ];

        this.setState({ chartOpts: chartOptions, showChart: true });
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
            <BRCharts opts={this.state.chartOpts} />
          </div>
        ) : (
          <div>Chart not available</div>
        )}
      </>
    );
  }
}
