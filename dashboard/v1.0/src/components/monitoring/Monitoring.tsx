import React from 'react';
import BRConnect from '../../utils/connection';
import {
  BRChartOpts,
  BRCharts,
  chartDefaultOptValues
} from '../layouts/BRCharts';
import Submenu from '../layouts/Submenu';

interface MonitoringModulePropsTypes {}

interface MonitoringModuleStateTypes {
  routes: object;
  sAddress: string;
  chartOpts: BRChartOpts[];
  showChart: boolean;  
}

export default class Monitoring extends React.Component<
  MonitoringModulePropsTypes,
  MonitoringModuleStateTypes
> {
  private connection: BRConnect;
  constructor(props: MonitoringModulePropsTypes) {
    super(props);

    this.connection = new BRConnect();
    const tmp: BRChartOpts[] = [chartDefaultOptValues];
    
    this.state = {
      // submenu address
      chartOpts: tmp,
      routes: {},
      sAddress: '',
      showChart: false
    };
  }

  public getAddressSubmenu = (sAddressParam: string) => {
    console.warn('addressSubmenu ', sAddressParam);
    this.setState({ sAddress: sAddressParam, showChart: false });
    this.connection.signalRequestResponseRouteFetchAllTimeSeries(sAddressParam)
    .then((res:any) => {
      const data: any[] = JSON.parse(res.data);
      const delay: BRChartOpts[] = [];
      const norTime: number[] = [];
      const timeStamp: string[] = [];

      let inst;
      for (inst of data) {
        delay.push(inst.delay);
        norTime.push(inst.relative);
        timeStamp.push(inst.timestamp);
      }
      const chartOptions: BRChartOpts[] = [
        {
          backgroundColor: 'rgba(75,192,192,0.4)',
          borderColor: 'rgba(75,192,192,1)',
          fill: false,
          label: 'Monitoring',
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
          yAxisValues: delay
        }
      ];

      this.setState({ chartOpts: chartOptions, showChart: true });
    })
    .catch((err: any) => {
      console.log("Error");
      console.log(err);
    });
  };

  public opts = (operation: string) => {
    switch (operation) {
      case 'start':
        this.connection.signalRequestResponseMonitoringStart().then(res => {
          if (res.data) {
            alert('Monitoring routine started');
          }
        });
        break;
      case 'stop':
        this.connection.signalRequestResponseMonitoringStop().then(res => {
          if (res.data) {
            alert('Monitoring routine stopped');
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
          onClick={() => this.opts('start')}>
            Start
          </button>

          <button className="button-operations btn btn-danger"
          onClick={() => this.opts('stop')}>
            Stop
          </button>
        </div>
        <div>
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
        </div>
      </>
    );
  }
}
