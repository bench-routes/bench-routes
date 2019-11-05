import React, { Component } from 'react';
import Submenu from '../layouts/Submenu';

interface JitterModulePropsTypes {}

interface JitterModuleStateTypes {
  routes: object;
  sAddress: string;
}

export default class Jitter extends Component<
  JitterModulePropsTypes,
  JitterModuleStateTypes
> {
  constructor(props: JitterModulePropsTypes) {
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
        <div>This is the jitter page.</div>
      </>
    );
  }
}
