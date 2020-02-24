import React, { FC } from 'react';
import { RouteComponentProps } from 'react-router-dom';
import { ServicesState } from './ServiceState';
import { RoutesSummary } from './RoutesSummary';

const Card = (head: string, component: JSX.Element) => (
  <div
    className="col-md-6 row"
    style={{ marginLeft: '0px', marginRight: '0px' }}
  >
    <div
      style={{
        margin: '1%',
        border: '1px solid #f1f1f3',
        borderRadius: '5px',
        width: '100%',
        overflowY: 'scroll',
        overflowX: 'hidden'
      }}
    >
      {head ? (
        <div
          style={{
            padding: '0% 2% 2% 2%',
            fontWeight: 'bold',
            borderBottom: '1px solid #f1f1f3',
            width: '100%'
          }}
        >
          {head}
        </div>
      ) : null}

      {component}
    </div>
  </div>
);

const Dashboard: FC<RouteComponentProps> = () => {
  return (
    <div className="row" style={{ margin: '3%' }}>
      {Card('Services state', <ServicesState />)}
      {Card('', <RoutesSummary />)}
    </div>
  );
};

export default Dashboard;
