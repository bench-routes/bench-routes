import React, { FC } from 'react';
import { useFetch } from '../../utils/useFetch';
import { HOST_IP } from '../../utils/types';
import { Alert, Spinner, Badge } from 'reactstrap';

interface RoutesSummaryResponses {
  testServicesRoutes: string[];
  monitoringRoutes: string[];
}

const section = (head: string, sub: string, labels: string[]) => (
  <div>
    <div style={{ paddingBottom: '1%', borderBottom: '1px solid #f1f1f3' }}>
      <strong>{head}</strong> &nbsp;
      {sub ? <sub>{sub}</sub> : null}
    </div>
    {labels.map((l: string, i: number) => (
      <div key={i}>
        <Badge color="primary">{l}</Badge>
      </div>
    ))}
  </div>
);

export const RoutesSummary: FC<{}> = () => {
  const { response, error } = useFetch<RoutesSummaryResponses>(
    `${HOST_IP}/routes-summary`
  );
  console.warn(response);

  if (error) {
    console.log(error);
    return (
      <Alert color="danger">
        Error: unable to fetch routes-summary details.
      </Alert>
    );
  } else if (response.data) {
    const routes = response.data;
    return (
      <div style={{ padding: '4%', height: '15vh' }}>
        {section(
          'Services',
          'ping, jitter, floodping',
          routes.testServicesRoutes
        )}
        {section('Monitoring', '', routes.monitoringRoutes)}
      </div>
    );
  }

  return (
    <Alert color="warning">
      Fetching... <Spinner color="info" />{' '}
    </Alert>
  );
};
