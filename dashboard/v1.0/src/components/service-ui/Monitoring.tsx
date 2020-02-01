import React from 'react';
import BRConnect from '../../utils/connection';
import Submenu from '../layouts/Submenu';
import { ChartOptions, Charts, ChartValues } from '../layouts/Charts';

interface MonitoringModulePropsTypes {}

interface MonitoringModuleStateTypes {
  routes: object;
  sAddress: string;
  showChart: boolean;
  options: ChartOptions[];
  responseLengthOpts: ChartOptions[];
}

export default class Monitoring extends React.Component<
  MonitoringModulePropsTypes,
  MonitoringModuleStateTypes
> {
  private connection: BRConnect;
  constructor(props: MonitoringModulePropsTypes) {
    super(props);

    this.connection = new BRConnect();
    const tmp: ChartOptions[] = [ChartValues()];
    this.state = {
      // submenu address
      routes: {},
      sAddress: '',
      showChart: false,
      options: tmp,
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

        const options: ChartOptions[] = [
          ChartValues(d, r, 'Delay', 'rgba(75,192,192,0.4)')
        ];

        const optionsResponse: ChartOptions[] = [
          ChartValues(d, resLength, 'Response-length', 'rgba(5,192,19,0.4)')
        ];

        this.setState({
          options,
          responseLengthOpts: optionsResponse,
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
            <Charts opts={this.state.options} />
            <Charts opts={this.state.responseLengthOpts} />
          </div>
        ) : (
          <div>Chart not available</div>
        )}
      </>
    );
  }
}
