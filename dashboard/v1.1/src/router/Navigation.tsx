import React, { FC } from 'react';
import { HashRouter as Router, Route, Switch } from 'react-router-dom';
import Dashboard from '../pages/Dashboard/Dashboard';
import FloodPing from '../pages/FloodPing';
import Jitter from '../pages/Jitter';
import Monitoring from '../pages/Monitoring';
import Ping from '../pages/Ping';
import Settings from '../pages/Settings';

interface NavigatorProps {
  updateLoader(status: boolean): void;
}

const Navigator: FC<NavigatorProps> = ({ updateLoader }) => {
  return (
    <Router>
      <Switch>
        <Route
          exact={true}
          path="/"
          render={props => <Dashboard updateLoader={updateLoader} />}
        />
        <Route path="/monitoring" component={Monitoring} />
        <Route path="/ping" component={Ping} />
        <Route path="/floodping" component={FloodPing} />
        <Route path="/jitter" component={Jitter} />
        <Route path="/settings" component={Settings} />
      </Switch>
    </Router>
  );
};

export default Navigator;
