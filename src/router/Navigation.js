import React from 'react';
import { HashRouter as Router, Switch, Route } from 'react-router-dom';
import Dashboard from '../components/dashboard/Dashboard';
import Monitoring from '../components/monitoring/Monitoring';
import Benchmarks from '../components/benchmarks/Benchmarks';

export default class Navigator extends React.Component {
  render() {
    return (
      <Router>
        <Switch>
          <Route exact path='/' component={Dashboard} />
          <Route path='/monitoring' component={Monitoring} />
          <Route path='/benchmarks' component={Benchmarks} />
        </Switch>
      </Router>
    );
  }
}
