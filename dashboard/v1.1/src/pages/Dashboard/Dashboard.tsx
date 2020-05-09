import React, { FC } from 'react';
import { RouteComponentProps } from 'react-router-dom';
import SystemMetrics from './SystemMetrics';
import Chart from 'react-apexcharts';

const Dashboard: FC<RouteComponentProps> = () => {
  return (
    <>
      BenchRoute Dashboard
      <SystemMetrics />
    </>
  );
};

export default Dashboard;
