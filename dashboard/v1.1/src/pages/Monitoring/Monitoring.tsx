import React, { FC, useEffect, useState } from 'react';
import { useFetch } from '../../utils/useFetch';
import { HOST_IP } from '../../utils/types';
import Matrix, { TimeSeriesPath, RouteDetails } from './Matrix';

import { Card, CardContent } from '@material-ui/core';
import Alert from '@material-ui/lab/Alert';

interface MonitoringProps {
  updateLoader(status: boolean): void;
}

const Monitoring: FC<MonitoringProps> = ({ updateLoader }) => {
  const { response, error } = useFetch<TimeSeriesPath[]>(
    `${HOST_IP}/get-route-time-series`
  );
  const [showRouteDetails, setShowRouteDetails] = useState<boolean>(false);
  const showDetails = (status: boolean, details: RouteDetails): void => {
    setShowRouteDetails(status);
  };
  useEffect(() => {
    updateLoader(true);
  }, [updateLoader]);
  if (error) {
    return (
      <Card>
        <CardContent>
          <h4>Monitoring</h4>
          <hr />
          <Alert severity="error">Unable to reach the service: error</Alert>
        </CardContent>
      </Card>
    );
  }
  if (!response.data) {
    return (
      <Card>
        <CardContent>
          <h4>Monitoring</h4>
          <hr />
          <Alert severity="info">Fetching data from sources</Alert>
        </CardContent>
      </Card>
    );
  }
  updateLoader(false);
  return (
    <Card>
      <CardContent>
        {!showRouteDetails ? (
          <>
            <h4>Monitoring</h4>
            <hr />
            <Matrix
              timeSeriesPath={response.data}
              showRouteDetails={showDetails}
            />
          </>
        ) : (
          <p>routeDetails</p>
        )}
      </CardContent>
    </Card>
  );
};

export default Monitoring;
