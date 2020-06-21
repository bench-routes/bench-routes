import React, { FC, useEffect, useState } from 'react';
import { useFetch } from '../../utils/useFetch';
import { HOST_IP } from '../../utils/types';
import Matrix from './Matrix';
import { TimeSeriesPath, RouteDetails } from '../../utils/queryTypes';
import RouteDetailsComponent from './RouteDetails';
import { Card, CardContent } from '@material-ui/core';
import Alert from '@material-ui/lab/Alert';
import Switch from '@material-ui/core/Switch';

interface MonitoringProps {
  updateLoader(status: boolean): void;
}

const ServicesState: FC<{}> = () => {
  const [isActive, setIsActive] = useState<boolean>(false);
  useEffect(() => {
    fetchState();
  }, []);
  const fetchState = () => {
    fetch(`${HOST_IP}/get-monitoring-services-state`)
      .then(res => res.json())
      .then((response: { status: string; data: string }) => {
        if (response.data === 'active') {
          setIsActive(true);
        } else {
          setIsActive(false);
        }
      });
  };
  const updateServicesState = () => {
    fetch(
      `${HOST_IP}/update-monitoring-services-state?state=${
        isActive ? 'stop' : 'start'
      }`
    )
      .then(res => res.json())
      .then((response: boolean) => {
        if (response) {
          setIsActive(!isActive);
        }
      });
  };
  return (
    <Switch
      checked={isActive}
      color="primary"
      onChange={() => updateServicesState()}
    />
  );
};

const Monitoring: FC<MonitoringProps> = ({ updateLoader }) => {
  const { response, error } = useFetch<TimeSeriesPath[]>(
    `${HOST_IP}/get-route-time-series`
  );
  const [showRouteDetails, setShowRouteDetails] = useState<boolean>(false);
  const [routeDetailsData, setRouteDetailsData] = useState<RouteDetails>();
  const showDetails = (status: boolean, details: RouteDetails): void => {
    setShowRouteDetails(status);
    setRouteDetailsData(details);
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
        {!showRouteDetails || !routeDetailsData ? (
          <>
            <h4>
              Monitoring <ServicesState />
            </h4>
            <hr />
            <Matrix
              timeSeriesPath={response.data}
              showRouteDetails={showDetails}
            />
          </>
        ) : (
          <RouteDetailsComponent
            routesChains={routeDetailsData}
            showDetails={showDetails}
          />
        )}
      </CardContent>
    </Card>
  );
};

export default Monitoring;
