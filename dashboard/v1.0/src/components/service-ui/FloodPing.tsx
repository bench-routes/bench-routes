import React, { FC, useState } from 'react';
import BRConnect from '../../utils/connection';
import { ChartOptions, Charts, ChartValues } from '../layouts/Charts';
import Submenu from '../layouts/Submenu';
import { Alert } from 'reactstrap';
import { opts } from './publicOpts';
import { getChartOptions } from './getChartOptions';

const FloodPing: FC<{}> = () => {
  const [chart, setChart] = useState({
    options: [ChartValues()],
    optionsPacketLoss: [ChartValues()],
    show: false
  });
  const connection = new BRConnect();

  const updateAddressSubmenu = (sAddressParam: string): void => {
    setChart({
      options: [ChartValues()],
      show: true,
      optionsPacketLoss: [ChartValues()]
    });
    connection
      .signalFloodPingRouteFetchAllTimeSeries(sAddressParam)
      .then((res: any) => {
        const data: any[] = JSON.parse(res.data);
        const norTime: number[] = [];
        const packetLoss: number[] = [];

        const options = getChartOptions(data);

        const optionsPacketLoss: ChartOptions[] = [
          ChartValues(
            norTime,
            packetLoss,
            'Packet-loss',
            'rgba(75,192,192,0.4)'
          )
        ];

        setChart({ options, show: true, optionsPacketLoss });
      });
  };

  return (
    <>
      <div className="btn-layout">
        <button
          className="button-operations btn btn-success"
          onClick={() => opts('start', connection, 'floodPing')}
        >
          Start
        </button>
        <button
          className="button-operations btn btn-danger"
          onClick={() => opts('stop', connection, 'floodPing')}
        >
          Stop
        </button>
      </div>
      <Submenu module="ping" submodule="" getAddress={updateAddressSubmenu} />
      {chart.show ? (
        <div style={{ overflow: 'scroll', height: '45%' }}>
          <Charts opts={chart.options} />
          <Charts opts={chart.optionsPacketLoss} />
        </div>
      ) : (
        <Alert color="warning">
          Select an option from the drop-down list for visualization.
        </Alert>
      )}
    </>
  );
};

export default FloodPing;
