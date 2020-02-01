import React, { FC } from 'react';
import { service_states, HOST_IP } from '../../utils/types';
import { useFetch } from '../../utils/useFetch';
import { Alert, Spinner, Badge } from 'reactstrap';

interface ConditionalBadgeProps {
  Key: string;
  value: string;
}

const ConditionalBadge: FC<ConditionalBadgeProps> = ({ Key, value }) => {
  if (value === 'active') {
    return <Badge color="warning">{`${Key}: ${value}`}</Badge>;
  } else {
    return <Badge color="danger">{`${Key}: ${value}`}</Badge>;
  }
};

export const ServicesState: FC<{}> = () => {
  const { response, error } = useFetch<service_states>(
    `${HOST_IP}/service-state`
  );

  if (error) {
    return <Alert type="error">Error: unable to reach the service.</Alert>;
  } else if (response.data) {
    const states: service_states = response.data;

    return (
      <div className="row" style={{ padding: '4%', height: '15vh' }}>
        <div className="col-md-6">
          <ConditionalBadge Key="Ping" value={states.ping} />
        </div>
        <div className="col-md-6">
          <ConditionalBadge Key="Jitter" value={states.jitter} />
        </div>
        <div className="col-md-6">
          <ConditionalBadge Key="Flood-Ping" value={states.floodping} />
        </div>
        <div className="col-md-6">
          <ConditionalBadge Key="Moitoring" value={states.monitoring} />
        </div>
      </div>
    );
  }

  return <Spinner color="info" />;
};
