import React from 'react';
import { HashRouter as Router, Route, Switch } from 'react-router-dom';
import FloodPing from '../components/benchmarks/FloodPing';
import Jitter from '../components/benchmarks/Jitter';
import Ping from '../components/benchmarks/Ping';
import Dashboard from '../components/dashboard/Dashboard';
import Monitoring from '../components/monitoring/Monitoring';
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
