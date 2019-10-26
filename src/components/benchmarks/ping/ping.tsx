import React from 'react';
import Submenu from '../../layouts/Submenu';

interface PingModulePropsTypes {}

interface PingModuleStateTypes {
  sAddress: string
}

export default class PingModule extends React.Component<
  PingModulePropsTypes,
  PingModuleStateTypes
> {
  constructor(props: PingModulePropsTypes){
    super(props);

    this.state = {
      // submenu ddress
      sAddress: ''
    }
  }

  getAddressSubmenu = (sAddressParam: string) => {
    this.setState({sAddress: sAddressParam})
  }

  render() {
    return (
      <Submenu module='ping' submodule='' getAddress={this.getAddressSubmenu} />
    );
  }
}