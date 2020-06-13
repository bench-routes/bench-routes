import React, { FC } from 'react';
import { HashRouter as Router, Route, Switch } from 'react-router-dom';
import Dashboard from '../pages/Dashboard/Dashboard';
import FloodPing from '../pages/FloodPing';
import JitterModule from '../pages/Jitter/JitterModule';
import Monitoring from '../pages/Monitoring/Monitoring';
import PingModule from '../pages/Ping/PingModule';
import Settings from '../pages/Settings';
import Input from '../pages/Input/Input';

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
        <Route
          exact={true}
          path="/monitoring"
          render={props => <Monitoring updateLoader={updateLoader} />}
        />
        <Route path="/ping" component={PingModule} />
        <Route path="/floodping" component={FloodPing} />
        <Route path="/jitter" component={JitterModule} />
        <Route path="/settings" component={Settings} />
        <Route path="/quick-input" component={Input} />
      </Switch>
    </Router>
  );
};

export default Navigator;
