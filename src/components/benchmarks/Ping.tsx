import React from 'react';
import Submenu from '../layouts/Submenu';

interface PingModulePropsTypes {}

interface PingModuleStateTypes {
  routes: object;
  sAddress: string;
}

export default class PingModule extends React.Component<
  PingModulePropsTypes,
  PingModuleStateTypes
> {
  constructor(props: PingModulePropsTypes) {
    super(props);

    this.state = {
      // submenu address
      routes: {},
      sAddress: ''
    };
  }

  public getAddressSubmenu = (sAddressParam: string) => {
    this.setState({ sAddress: sAddressParam });
  };

  public render() {
    return (
      <>
        <Submenu
          module="ping"
          submodule=""
          getAddress={this.getAddressSubmenu}
        />
        <div>This is the ping page.</div>
      </>
    );
  }
}
