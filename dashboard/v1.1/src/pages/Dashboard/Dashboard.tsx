import React, { FC } from 'react';
import SystemMetrics from './SystemMetrics';
import JournalMetrics from './JournalMetrics';
import { Card, CardContent } from '@material-ui/core';
import GraphWrapper from '../../layouts/GraphWrapper';

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
        <GraphWrapper
          child={<SystemMetrics showLoader={updateLoader} />}
          isMonitoring={false}
        />
        <hr />
        <GraphWrapper child={<JournalMetrics />} isMonitoring={false} />
      </CardContent>
    </Card>
  );
};

export default Dashboard;
