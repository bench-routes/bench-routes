import React, { FC, useEffect, useState, useCallback } from 'react';
import { useFetch } from '../../utils/useFetch';
import { HOST_IP } from '../../utils/types';
import Matrix from './Matrix';
import { TimeSeriesPath, RouteDetails } from '../../utils/queryTypes';
import RouteDetailsComponent from './RouteDetails';
import { Card, CardContent, Tooltip, Fab, makeStyles } from '@material-ui/core';
import Alert from '@material-ui/lab/Alert';
import Switch from '@material-ui/core/Switch';
import { PostAdd as PostAddIcon } from '@material-ui/icons';
import { Link } from 'react-router-dom';

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
      <>
        <Card>
          <CardContent>
            <h4>Monitoring</h4>
            <hr />
            <Alert severity="error">Unable to reach the service: error</Alert>
          </CardContent>
        </Card>
        <InputFab />
      </>
    );
  }
  if (!response.data) {
    return (
      <>
        <Card>
          <CardContent>
            <h4>Monitoring</h4>
            <hr />
            <Alert severity="info">Fetching data from sources</Alert>
          </CardContent>
        </Card>
        <InputFab />
      </>
    );
  }
  updateLoader(false);
  return (
    <>
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
              <hr />
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
      <InputFab />
    </>
  );
};

export default Monitoring;

const useStyles = makeStyles(theme => ({
  inputFab: {
    position: 'fixed',
    right: '2rem',
    bottom: '2rem'
  }
}));

//Floating Action Button for Quick Input
const InputFab = () => {
  const classes = useStyles();
  return (
    <Link to="/quick-input">
      <Tooltip placement="top" title="Quick Input">
        <Fab
          className={classes.inputFab}
          color="primary"
          aria-label="quick-input"
        >
          <PostAddIcon />
        </Fab>
      </Tooltip>
    </Link>
  );
};
