import React from 'react';
import { HashRouter as Router, Route, Switch } from 'react-router-dom';
import Dashboard from '../components/dashboard/Dashboard';
import FloodPing from '../components/service-ui/FloodPing';
import Jitter from '../components/service-ui/Jitter';
import Monitoring from '../components/service-ui/Monitoring';
import Ping from '../components/service-ui/Ping';
import Settings from '../components/settings/Settings';

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
