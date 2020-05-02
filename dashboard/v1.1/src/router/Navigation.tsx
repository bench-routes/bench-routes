import React from 'react';
import { HashRouter as Router, Route, Switch } from 'react-router-dom';
import Dashboard from '../pages/Dashboard/Dashboard';
import FloodPing from '../pages/FloodPing';
import Jitter from '../pages/Jitter';
import Monitoring from '../pages/Monitoring';
import Ping from '../pages/Ping';
import Settings from '../pages/Settings';

const Navigator = () => (
  <Router>
    <Switch>
      <Route exact={true} path="/" component={Dashboard} />
      <Route path="/monitoring" component={Monitoring} />
      <Route path="/ping" component={Ping} />
      <Route path="/floodping" component={FloodPing} />
      <Route path="/jitter" component={Jitter} />
      <Route path="/settings" component={Settings} />
    </Switch>
  </Router>
);

export default Navigator;
