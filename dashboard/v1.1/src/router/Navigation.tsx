import React, { FC } from 'react';
import { HashRouter as Router, Route, Switch } from 'react-router-dom';
import QuickInputFab from '../layouts/QuickInputFab';
import Config from '../pages/Config/Config';
import Dashboard from '../pages/Dashboard/Dashboard';
import FloodPing from '../pages/FloodPing';
import Input from '../pages/Input/Input';
import JitterModule from '../pages/Jitter/JitterModule';
import Monitoring from '../pages/Monitoring/Monitoring';
import PingModule from '../pages/Ping/PingModule';

interface NavigatorProps {
  updateLoader(status: boolean): void;
}

const Navigator: FC<NavigatorProps> = ({ updateLoader }) => {
  return (
    <Router>
      {/* Floating Action Button for Quick Route Input */}
      <QuickInputFab />
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
        <Route path="/quick-input" component={Input} />
        <Route path="/configurations" component={Config} />
      </Switch>
    </Router>
  );
};

export default Navigator;
