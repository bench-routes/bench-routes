import React, { FC, useState } from 'react';
import { RouteComponentProps } from 'react-router-dom';
import SystemMetrics from './SystemMetrics';
import { Card, CardContent } from '@material-ui/core';
import LinearProgress from '@material-ui/core/LinearProgress';

const Dashboard: FC<RouteComponentProps> = () => {
  const [systemMetricsDone, setSystemMetricsDone] = useState<boolean>(false);
  const systemMetricsLoaded = (status: boolean) => {
    if (status) {
      setSystemMetricsDone(true);
    } else {
      setSystemMetricsDone(false);
    }
  };
  return (
    <Card>
      <CardContent>
        <h4>Dashboard</h4>
        {!systemMetricsDone ? <LinearProgress /> : <hr />}
        <SystemMetrics done={systemMetricsLoaded} />
      </CardContent>
    </Card>
  );
};

export default Dashboard;
