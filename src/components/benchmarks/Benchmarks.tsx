import React from 'react';
import Submenu from '../layouts/Submenu';

export default class Benchmarks extends React.Component {
  state = {
    addressSubmenu: ''
  }

  getAddressSubmenu = (s: string) => {
    console.log('here')
    this.setState({addressSubmenu: s});
    console.log(this.state.addressSubmenu)
  }

  render() {
    return (
      <>
        <Submenu module='benchmark' submodule='' getAddress={this.getAddressSubmenu} />
        This is Benchmarking
      </>
    );
  }
}
