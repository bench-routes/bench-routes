import React from 'react';
import BRConnect from '../../utils/connection';
import Submenu from '../layouts/Submenu';

interface MonitoringModulePropsTypes {}

interface MonitoringModuleStateTypes {
  routes: object;
  sAddress: string;
}

export default class Monitoring extends React.Component<
  MonitoringModulePropsTypes,
  MonitoringModuleStateTypes
> {
  private connection: BRConnect;
  constructor(props: MonitoringModulePropsTypes) {
    super(props);

    this.connection = new BRConnect();
    this.state = {
      // submenu address
      routes: {},
      sAddress: ''
    };
  }

  public getAddressSubmenu = (sAddressParam: string) => {
    // Note that teh sAddressParam should not contain any "/" in the end. If exists, trim it.
    sAddressParam = sAddressParam.substring(0, sAddressParam.length - 3);
    let arr: string[] = sAddressParam.split(' ');
    sAddressParam = sAddressParam.split(' ')[2];
    console.warn('addressSubmenu ', sAddressParam);
    this.setState({ sAddress: sAddressParam });
    this.connection.signalRequestResponseRouteFetchAllTimeSeries(sAddressParam);
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
        <div>This is the ping page.</div>
      </>
    );
  }
}
