import React, { FC } from 'react';
import { RouteComponentProps } from 'react-router-dom';
import { PersistentConnection } from '../Services/socket';

const Dashboard: FC<RouteComponentProps> = () => {
  const conn = new PersistentConnection().getSocketInstance();
  conn.onopen = () => {
    conn.send('connected');
  };
  return <div>BenchRoute Dashboard</div>;
};

export default Dashboard;
