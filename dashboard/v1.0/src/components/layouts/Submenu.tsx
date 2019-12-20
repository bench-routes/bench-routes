import React from 'react';
import BRConnect from '../../utils/connection';

interface SubmenuPropsTypes {
  module: string;
  submodule: string;
  getAddress: (sAddress: string) => void;
}

interface RouterType {
  headers: object[];
  method: string;
  params: object[];
  url: string;
}

interface SubmenuStateTypes {
  module: string;
  urlSlot: string;
  methodSlot: string;
  routeSlot: string;
  paramSlot: string;
  routes: any;
}

export default class Submenu extends React.Component<
  SubmenuPropsTypes,
  SubmenuStateTypes
> {
  private BRinstance: BRConnect;
  constructor(props: SubmenuPropsTypes) {
    super(props);

    const tmp: RouterType[] = [];
    this.BRinstance = new BRConnect();
    this.state = {
      methodSlot: '',
      module: '',
      paramSlot: '',
      routeSlot: '',
      routes: tmp,
      urlSlot: ''
    };
  }

  public componentDidMount() {
    this.setState({ module: this.props.module });
    this.BRinstance.routeDetails().then((res: any) => {
      this.setState({ routes: res.routes });
    });
  }

  public sendAddress = () => {
    switch (this.props.module) {
      case 'ping':
        this.props.getAddress(this.state.urlSlot);
        break;
      case 'jitter':
        this.props.getAddress(this.state.urlSlot);
        break;
      default:
        this.props.getAddress(
          this.state.urlSlot +
            '/' +
            this.state.methodSlot +
            '/' +
            this.state.routeSlot +
            '/' +
            this.state.paramSlot
        );
    }
  };

  public render() {
    return (
      <>
        {/* submenu container */}
        <div className="submenu-container">
          {/* Drop-down components */}
          <div>
            <span>
              {/* IP / Domain */}
              {(this.props.module === 'ping' ||
                this.props.module === 'jitter' ||
                this.props.module === 'flood-ping') &&
              this.props.submodule.length === 0 ? (
                <span>
                  <select
                    className="submenu-style-general"
                    onChange={e => this.setState({ urlSlot: e.target.value })}
                  >
                    <option />
                    {this.state.routes.length !== 0
                      ? this.state.routes.map((val: RouterType, id: number) => (
                          <option key={id} value={val.url}>
                            {val.url}
                          </option>
                        ))
                      : null}
                  </select>
                </span>
              ) : (
                <span>
                  <select
                    className="submenu-style-general"
                    onChange={e => this.setState({ urlSlot: e.target.value })}
                  >
                    <option />
                    {this.state.routes.length !== 0
                      ? this.state.routes.map((val: any, id: number) => (
                          <option
                            key={id}
                            value={
                              val.method + '  ' + val.url + '/' + val.route
                            }
                          >
                            {val.method + '  ' + val.url + '/' + val.route}
                          </option>
                        ))
                      : null}
                  </select>
                </span>
              )}
            </span>
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
