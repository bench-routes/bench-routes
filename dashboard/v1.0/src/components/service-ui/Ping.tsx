import React, { FC, useState } from 'react';
import BRConnect from '../../utils/connection';
import { ChartOptions, Charts, ChartValues } from '../layouts/Charts';
import Submenu from '../layouts/Submenu';
import { Alert } from 'reactstrap';
import {opts} from './opts';

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
        const yMin: ChartOptions[] = [];
        const yMean: ChartOptions[] = [];
        const yMax: ChartOptions[] = [];
        const yMdev: ChartOptions[] = [];
        const norTime: number[] = [];
        const timeStamp: string[] = [];

        if (data.length === 0) {
          // Probably send the required information
          // to the user via br-logger
          console.log('No data from the url');
        } else {
          let inst;
          for (inst of data) {
            yMin.push(inst.Min);
            yMean.push(inst.Mean);
            yMax.push(inst.Max);
            yMdev.push(inst.Mdev);
            norTime.push(inst.relative);
            timeStamp.push(inst.timestamp);
          }
        }

        const options: ChartOptions[] = [
          ChartValues(norTime, yMin, 'Minimum', 'rgba(75,192,192,0.4)'),
          ChartValues(norTime, yMean, 'Mean', 'rgba(75,192,2,0.4)'),
          ChartValues(norTime, yMax, 'Maximum', 'rgba(5,192,19,0.4)'),
          ChartValues(norTime, yMdev, 'Standard-Deviation', 'rgba(7,12,19,0.4)')
        ];

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
