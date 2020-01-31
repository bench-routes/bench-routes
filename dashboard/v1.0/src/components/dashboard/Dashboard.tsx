import React, { FC } from 'react';
import { RouteComponentProps } from 'react-router-dom';
import { ServicesState } from './ServiceState';
import { MonitoringSummary } from './MonitoringSummary';

const Dashboard: FC<RouteComponentProps> = () => {
  return (
    <div className="row" style={{ padding: '3%' }}>
      <div
        className="col-md-6 row"
        style={{ border: '1px solid #f1f1f3', borderRadius: '5px' }}
      >
        <div
          style={{
            padding: '2%',
            fontWeight: 'bold',
            borderBottom: '1px solid #f1f1f3',
            width: '100%'
          }}
        >
          Services state
        </div>
        <ServicesState />
      </div>
      <MonitoringSummary />
    </div>
  );
};

export default Dashboard;
