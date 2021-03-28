import React, { FC, useEffect, useState, useCallback } from 'react';
import { useFetch } from '../../utils/useFetch';
import { HOST_IP } from '../../utils/types';
import Matrix from './Matrix';
import { TimeSeriesPath, RouteDetails } from '../../utils/queryTypes';
import RouteDetailsComponent from './RouteDetails';
import { Card, CardContent, Tooltip } from '@material-ui/core';
import Alert from '@material-ui/lab/Alert';
import Switch from '@material-ui/core/Switch';

interface MonitoringProps {
  updateLoader(status: boolean): void;
}

interface ServiceStateProps {
  active(status: boolean): void;
}

const ServicesState: FC<ServiceStateProps> = ({ active }) => {
  const [isActive, setIsActive] = useState<boolean>(false);

  const fetchState = useCallback(async () => {
    const raw = await fetch(`${HOST_IP}/get-monitoring-services-state`);
    const JSON = (await raw.json()) as { status: string; data: string };
    if (JSON.data === 'active') {
      setIsActive(true);
      active(true);
    } else {
      setIsActive(false);
      active(false);
    }
  }, [active]);

  useEffect(() => {
    fetchState();
  }, [isActive, fetchState]);

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
          active(!isActive);
        }
      });
  };

  return (
    <Tooltip
      title={(isActive ? 'Stop' : 'Start') + ' Monitoring'}
      aria-label={(isActive ? 'Stop' : 'Start') + ' Monitoring'}
    >
      <Switch
        checked={isActive}
        color="primary"
        onChange={() => updateServicesState()}
      />
    </Tooltip>
  );
};

const Monitoring: FC<MonitoringProps> = ({ updateLoader }) => {
  const { response, error } = useFetch<TimeSeriesPath[]>(
    `${HOST_IP}/get-route-time-series`
  );
  const [showRouteDetails, setShowRouteDetails] = useState<boolean>(false);
  const [routeDetailsData, setRouteDetailsData] = useState<RouteDetails>();
  const [isMonitoringActive, setIsMonitoringActive] = useState<boolean>(false);
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
              Monitoring{' '}
              <ServicesState
                active={(status: boolean) => setIsMonitoringActive(status)}
              />
            </h4>
            <Matrix
              timeSeriesPath={response.data}
              showRouteDetails={showDetails}
              isMonitoringActive={isMonitoringActive}
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
