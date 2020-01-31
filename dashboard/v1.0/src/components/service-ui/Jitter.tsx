import React, { FC, useState } from 'react';
import BRConnect from '../../utils/connection';
import { ChartOptions, Charts, ChartValues } from '../layouts/Charts';
import Submenu from '../layouts/Submenu';
import { Alert } from 'reactstrap';

const Jitter: FC<{}> = () => {
  const [chart, setChart] = useState({
    options: [ChartValues()],
    show: false
  });

  const connection = new BRConnect();

  const opts = (operation: string) => {
    switch (operation) {
      case 'start':
        connection.signalJitterStart().then(res => {
          if (res.data) {
            alert('Jitter routine started');
          }
        });
        break;
      case 'stop':
        connection.signalJitterStop().then(res => {
          if (res.data) {
            alert('Jitter routine stopped');
          }
        });
        break;
    }
  };

  const getAddressSubmenu = (sAddressParam: string): void => {
    setChart({ options: [ChartValues()], show: false });
    connection
      .signalJitterRouteFetchAllTimeSeries(sAddressParam)
      .then((res: any) => {
        const data: any[] = JSON.parse(res.data);
        const jitter: ChartOptions[] = [];
        const norTime: number[] = [];
        const timeStamp: string[] = [];

        let inst;
        for (inst of data) {
          jitter.push(inst.datapoint);
          norTime.push(inst.relative);
          timeStamp.push(inst.timestamp);
        }

        const options: ChartOptions[] = [
          ChartValues(norTime, jitter, 'Jitter', 'rgba(75,192,192,0.4)')
        ];
        setChart({ options, show: true });
      });
  };

  return (
    <>
      <div className="btn-layout">
        <button
          className="button-operations btn btn-success"
          onClick={() => opts('start')}
        >
          Start
        </button>
        <button
          className="button-operations btn btn-danger"
          onClick={() => opts('stop')}
        >
          Stop
        </button>
      </div>
      <Submenu module="jitter" submodule="" getAddress={getAddressSubmenu} />
      {chart.show ? (
        <div>
          <Charts opts={chart.options} />
        </div>
      ) : (
        <Alert color="warning">
          Select an option from the drop-down list for visualization.
        </Alert>
      )}
    </>
  );
};

export default Jitter;
