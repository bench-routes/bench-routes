import React from 'react';

interface SubmenuPropsTypes {
  module: string;
  submodule: string;
  getAddress: (sAddress: string) => void;
}

interface SubmenuStateTypes {
  module: string;
  urlSlot: string;
  methodSlot: string;
  routeSlot: string;
  paramSlot: string;
}

export default class Submenu extends React.Component<
  SubmenuPropsTypes,
  SubmenuStateTypes
> {
  constructor(props: SubmenuPropsTypes) {
    super(props);

    this.state = {
      methodSlot: '',
      module: '',
      paramSlot: '',
      routeSlot: '',
      urlSlot: ''
    };
  }

  public sendAddress = () => {
    this.props.getAddress(
      this.state.urlSlot +
        '/' +
        this.state.methodSlot +
        '/' +
        this.state.routeSlot +
        '/' +
        this.state.paramSlot
    );
  };

  public componentDidMount() {
    this.setState({ module: this.props.module });
  }

  public render() {
    return (
      <>
        {/* submenu container */}
        <div className="submenu-container">
          {/* Drop-down components */}
          <div>
            <div>
              {/* IP / Domain */}
              {this.props.module === 'benchmark' &&
              this.props.submodule.length === 0 ? (
                <div>
                  <select
                    className="submenu-style-general"
                    onChange={e => this.setState({ urlSlot: e.target.value })}
                  >
                    <option></option>
                    <option>google.co.in</option>
                    <option>bing.com</option>
                    <option>yahoo.com</option>
                  </select>
                </div>
              ) : null}
            </div>
          </div>

          <div className="float-right">
            <button
              onClick={this.sendAddress}
              className="submenu-show-graph btn-primary"
            >
              Show
            </button>
          </div>
        </div>
      </>
    );
  }
}
