import React, { FC, useState } from 'react';
import BRConnect from '../../utils/connection';
import { Charts, ChartValues } from '../layouts/Charts';
import Submenu from '../layouts/Submenu';
import { Alert } from 'reactstrap';
import { opts } from './publicOpts';
import { getChartOptions } from './getChartOptions';
const Ping: FC<{}> = () => {
  const [chart, setChart] = useState({
    options: [ChartValues()],
    show: false
  });
  const connection: BRConnect = new BRConnect();

  const updateAddressSubmenu = (sAddressParam: string): void => {
    setChart({ options: [ChartValues()], show: false });
    connection
      .signalPingRouteFetchAllTimeSeries(sAddressParam)
      .then((res: any) => {
        let data: any[] = JSON.parse(res.data) || [];

        const options = getChartOptions(data);
        setChart({ options, show: true });
      });
  };

  return (
    <>
      <div className="btn-layout">
        <button
          className="button-operations btn btn-success"
          onClick={() => opts('start', connection, 'ping')}
        >
          Start
        </button>
        <button
          className="button-operations btn btn-danger"
          onClick={() => opts('stop', connection, 'ping')}
        >
          Stop
        </button>
      </div>
      <div>
        <Submenu module="ping" submodule="" getAddress={updateAddressSubmenu} />
        {chart.show ? (
          <div>
            <Charts opts={chart.options} />
          </div>
        ) : (
          <Alert color="warning">
            Select an option from the drop-down list for visualization.
          </Alert>
        )}
      </div>
    </>
  );
};

export default Ping;
