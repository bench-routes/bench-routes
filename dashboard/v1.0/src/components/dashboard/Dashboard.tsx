import React, { FC } from 'react';
import { RouteComponentProps, NavLink } from 'react-router-dom';

// export default class Dashboard extends React.Component<{}> {


//   public render() {
//     return <div>This is a Dashboard</div>;
//   }
// }

const Dashboard: FC<RouteComponentProps> = () => {
  fetch('http://localhost:9090/service-state')
    .then(res => res.json())
    .then(response => {
      console.warn('this is response', response)
    })

  return <div>This is a dashboard page</div>;
};

export default Dashboard;
