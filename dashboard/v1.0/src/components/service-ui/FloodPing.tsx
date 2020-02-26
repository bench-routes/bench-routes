import React, { FC, useState } from 'react';
import BRConnect from '../../utils/connection';
import { ChartOptions, Charts, ChartValues } from '../layouts/Charts';
import Submenu from '../layouts/Submenu';
import { Alert } from 'reactstrap';

const FloodPing: FC<{}> = () => {
  const [chart, setChart] = useState({
    options: [ChartValues()],
    optionsPacketLoss: [ChartValues()],
    show: false
  });
  const connection = new BRConnect();

  const opts = (operation: string) => {
    switch (operation) {
      case 'start':
        connection.signalFloodPingStart().then(res => {
          if (res.data) {
            alert('Flood Ping routine started');
          }
        });
        break;
      case 'stop':
        connection.signalFloodPingStop().then(res => {
          if (res.data) {
            alert('Flood Ping routine stopped');
          }
        });
        break;
    }
  };

  const setAddressSubmenu = (sAddressParam: string): void => {
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
        const timeStamp: string[] = [];
        const yMin: ChartOptions[] = [];
        const yMean: ChartOptions[] = [];
        const yMax: ChartOptions[] = [];
        const yMdev: ChartOptions[] = [];
        const packetLoss: number[] = [];

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
            packetLoss.push(inst.PacketLoss);
          }
        }

        const options: ChartOptions[] = [
          ChartValues(norTime, yMin, 'Minimum', 'rgba(75,192,192,0.4)'),
          ChartValues(norTime, yMean, 'Mean', 'rgba(75,192,2,0.4)'),
          ChartValues(norTime, yMax, 'Maximum', 'rgba(5,192,19,0.4)'),
          ChartValues(norTime, yMdev, 'Standard-Deviation', 'rgba(7,12,19,0.4)')
        ];

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
          onClick={() => opts('start')}
        >
          Start
        </button>
        <button
          className="button-operations btn btn-danger"
          onClick={() => opts('start')}
        >
          Stop
        </button>
      </div>
      <Submenu module="ping" submodule="" getAddress={setAddressSubmenu} />
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
