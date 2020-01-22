import React from 'react';
import BRConnect from '../../utils/connection';
import Submenu from '../layouts/Submenu';
import {
  BRChartOpts,
  BRCharts,
  chartDefaultOptValues
} from '../layouts/BRCharts';

interface MonitoringModulePropsTypes {}

interface MonitoringModuleStateTypes {
  routes: object;
  sAddress: string;
  showChart: boolean;
  chartOpts: BRChartOpts[];
  responseLengthOpts: BRChartOpts[];
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
      routes: {},
      sAddress: '',
      showChart: false,
      chartOpts: tmp,
      responseLengthOpts: tmp
    };
  }

  public getAddressSubmenu = (sAddressParam: string) => {
    // Note that the sAddressParam should not contain any "/" in the end. If exists, trim it.
    this.setState({ showChart: false });
    sAddressParam = sAddressParam.substring(0, sAddressParam.length - 3);
    sAddressParam = sAddressParam.split(' ')[2];
    console.warn('addressSubmenu ', sAddressParam);
    this.setState({ sAddress: sAddressParam });
    this.connection
      .signalRequestResponseRouteFetchAllTimeSeries(sAddressParam)
      .then((res: any) => {
        const data: any = JSON.parse(res.data);
        const d: any = [];
        const r: any = [];
        const resLength: any = [];
        for (const inst of data) {
          r.push(inst.delay);
          resLength.push(inst.resLength);
          d.push(inst.relative);
        }

        const chartOpts: BRChartOpts[] = [
          {
            backgroundColor: 'rgba(75,192,192,0.4)',
            borderColor: 'rgba(75,192,192,1)',
            fill: false,
            label: 'Delay',
            lineTension: 0.1,
            pBackgroundColor: '#fff',
            pBorderColor: 'rgba(75,192,2,1)',
            pBorderWidth: 1,
            pHoverBackgroundColor: 'rgba(75,192,192,1)',
            pHoverBorderColor: 'rgba(220,220,220,1)',
            pHoverBorderWidth: 2,
            pHoverRadius: 5,
            pRadius: 1,
            xAxisValues: d,
            yAxisValues: r
          }
        ];
        const resChartOpts: BRChartOpts[] = [
          {
            backgroundColor: 'rgba(5,192,19,0.4)',
            borderColor: 'rgba(5,192,19,0.4)',
            fill: false,
            label: 'Response-length',
            lineTension: 0.1,
            pBackgroundColor: '#fff',
            pBorderColor: 'rgba(75,192,2,1)',
            pBorderWidth: 1,
            pHoverBackgroundColor: 'rgba(75,192,192,1)',
            pHoverBorderColor: 'rgba(220,220,220,1)',
            pHoverBorderWidth: 2,
            pHoverRadius: 5,
            pRadius: 1,
            xAxisValues: d,
            yAxisValues: resLength
          }
        ];
        this.setState({
          chartOpts,
          responseLengthOpts: resChartOpts,
          showChart: true
        });
      });
  };

  public render() {
    return (
      <>
        <div className="btn-layout">
          {/* operations */}
          <button className="button-operations btn btn-success">Start</button>
          <button className="button-operations btn btn-danger">Stop</button>
        </div>
        <Submenu
          module="monitoring"
          submodule=""
          getAddress={this.getAddressSubmenu}
        />
        {this.state.showChart ? (
          <div style={{ overflowY: 'scroll', height: '45%' }}>
            <BRCharts opts={this.state.chartOpts} />
            <BRCharts opts={this.state.responseLengthOpts} />
          </div>
        ) : (
          <div>Chart not available</div>
        )}
      </>
    );
  }
}
