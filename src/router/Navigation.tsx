import React from 'react';
import { HashRouter as Router, Route, Switch } from 'react-router-dom';
import Benchmarks from '../components/benchmarks/Benchmarks';
import Dashboard from '../components/dashboard/Dashboard';
import Monitoring from '../components/monitoring/Monitoring';
import Settings from '../components/settings/Settings';

const Navigator = () => (
  <Router>
    <Switch>
      <Route exact path="/" component={Dashboard} />
      <Route path="/monitoring" component={Monitoring} />
      <Route path="/benchmarks" component={Benchmarks} />
      <Route path="/settings" component={Settings} />
    </Switch>
  </Router>
);

export default Navigator;
