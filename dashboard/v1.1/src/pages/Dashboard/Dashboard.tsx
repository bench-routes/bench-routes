import React, { FC } from 'react';
import SystemMetrics from './SystemMetrics';
import JournalMetrics from './JournalMetrics';
import { Card, CardContent } from '@material-ui/core';

interface DashboardProps {
  updateLoader(status: boolean): void;
  darkMode(status: boolean): void;
}

const Dashboard: FC<DashboardProps> = ({ updateLoader, darkMode }) => {
  updateLoader(true);
  return (
    <Card>
      <CardContent>
        <h4>Dashboard</h4>
        <hr />
        <SystemMetrics showLoader={updateLoader} darkMode={darkMode} />
        <hr />
        <JournalMetrics darkMode={darkMode} />
      </CardContent>
    </Card>
  );
};

export default Dashboard;
