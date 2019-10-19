import React from 'react';
import './layouts.style.css';
import { HashRouter as Router, Link } from 'react-router-dom';

export default class Sidebar extends React.Component {
    render() {
      return (
        <Router>
          <div className='sidebar'>
            <Link to='/' style={{ textDecoration: 'none' }}>
              <div>
                <div className='sidebar-inner'>
                  <img src='assets/icons/dashboard-icon.svg' className='sidebar-inner' alt='dashboard' />
                  <div className='sidebar-head sidebar-inner'>Dashboard</div>
                </div>
              </div>
            </Link>
            <Link to='/monitoring' style={{ textDecoration: 'none' }}>
              <div>
                <div className='sidebar-inner'>
                  <img src='assets/icons/monitoring-icon.svg' className='sidebar-inner' alt='monitoring' />
                  <div className='sidebar-head sidebar-inner'>Monitoring</div>
                </div>
              </div>
            </Link>
            <Link to='/benchmarks' style={{ textDecoration: 'none' }}>
              <div>
                <div className='sidebar-inner'>
                  <img src='assets/icons/bench-icon.svg' className='sidebar-inner' alt='Benchmarks' />
                  <div className='sidebar-head sidebar-inner'>Benchmarks</div>
                </div>
              </div>
            </Link>
            
          </div>
        </Router>
      );
    }
}