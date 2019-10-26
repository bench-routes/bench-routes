import React from 'react';

export default class Submenu extends React.Component<{
  module: string,
  submodule: string,
  getAddress: (sAddress: string) => void
}> {
  state = {
    module: '',
    urlSlot: '',
    methodSlot: '',
    routeSlot: '',
    paramSlot: ''
  }

  sendAddress = () => {
    this.props.getAddress(this.state.urlSlot + '/' + this.state.methodSlot + '/' + this.state.routeSlot + '/' + this.state.paramSlot);
  }

  componentDidMount() {
    this.setState({module: this.props.module})
  }

  render() {
    return (
      <>
        {/* submenu container */}
        <div className='submenu-container'>
          {/* Drop-down components */}
          <div>
            <div>
              {/* IP / Domain */}
              {
                this.props.module === 'benchmark' && this.props.submodule.length === 0
                ?
                <div>
                  <select className='submenu-style-general'
                    onChange={(e) => this.setState({urlSlot: e.target.value})}
                  >
                    <option></option>
                    <option>google.co.in</option>
                    <option>bing.com</option>
                    <option>yahoo.com</option>
                  </select>
                </div>
                : null
              }
            </div>
          </div>

          <div className='float-right'>
            <button
              onClick={this.sendAddress}
              className='submenu-show-graph btn-primary'
            >
              Show
            </button>
          </div>
        </div>
      </>
    );
  }
}