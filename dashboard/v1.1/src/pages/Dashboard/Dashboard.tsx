import React, { FC } from 'react';
import SystemMetrics from './SystemMetrics';
import { Card, CardContent } from '@material-ui/core';

interface DashboardProps {
  updateLoader(status: boolean): void;
}

const Dashboard: FC<DashboardProps> = ({ updateLoader }) => {
  updateLoader(true);
  return (
    <Card>
      <CardContent>
        <h4>Dashboard</h4>
        <SystemMetrics showLoader={updateLoader} />
      </CardContent>
    </Card>
  );
};

export default Dashboard;
