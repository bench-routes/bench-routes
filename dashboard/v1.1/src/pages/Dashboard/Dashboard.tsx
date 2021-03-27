import React, { FC } from 'react';
import SystemMetrics from './SystemMetrics';
import JournalMetrics from './JournalMetrics';
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
        <hr />
        <SystemMetrics showLoader={updateLoader} />
        <hr />
        <JournalMetrics />
      </CardContent>
    </Card>
  );
};

export default Dashboard;
