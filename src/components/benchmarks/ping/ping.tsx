import React from 'react';
import Submenu from '../../layouts/Submenu';

export default class PingModule extends React.Component {
  state = {
    // submenu ddress
    sAddress: ''
  }

  getAddressSubmenu = (sAddress: string) => {
    this.state.sAddress = sAddress;
  }

  render() {
    return (
      <Submenu module='ping' submodule='' getAddress={this.getAddressSubmenu} />
    );
  }
}