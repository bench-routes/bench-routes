import React from 'react';
import Submenu from '../layouts/Submenu';

interface BenchmarksPropsTypes {}

interface BenchmarksStateTypes {
  addressSubmenu: string
}

export default class Benchmarks extends React.Component<
  BenchmarksPropsTypes,
  BenchmarksStateTypes
> {
  constructor(props: BenchmarksPropsTypes) {
    super(props);

    this.state = {
      addressSubmenu: ''
    }
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
