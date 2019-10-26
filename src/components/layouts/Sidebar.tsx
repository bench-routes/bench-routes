import React from 'react';
import { HashRouter as Router, Link } from 'react-router-dom';
import './layouts.style.css';

export default class Sidebar extends React.Component<{}> {
  public render() {
    return (
      <Router>
        <div className="sidebar">
          <div className="sidebar-content">
            <Link to="/" style={{ textDecoration: 'none' }}>
              <div>
                <div className="sidebar-inner">
                  <img
                    src="assets/icons/dashboard-icon.svg"
                    className="sidebar-inner"
                    alt="dashboard"
                  />
                  <div className="sidebar-head sidebar-inner">Dashboard</div>
                </div>
              </div>
            </Link>
            <Link to="/monitoring" style={{ textDecoration: 'none' }}>
              <div>
                <div className="sidebar-inner">
                  <img
                    src="assets/icons/monitoring-icon.svg"
                    className="sidebar-inner"
                    alt="monitoring"
                  />
                  <div className="sidebar-head sidebar-inner">Monitoring</div>
                </div>
              </div>
            </Link>
            <Link to="/benchmarks" style={{ textDecoration: 'none' }}>
              <div>
                <div className="sidebar-inner">
                  <img
                    src="assets/icons/bench-icon.svg"
                    className="sidebar-inner"
                    alt="Benchmarks"
                  />
                  <div className="sidebar-head sidebar-inner">Benchmarks</div>
                </div>
              </div>
            </Link>
            <div className="sidebar-bottom-links">
              <Link to="/settings" style={{ textDecoration: 'none' }}>
                <div>
                  <div className="sidebar-inner">
                    <img
                      src="assets/icons/settings-icon.svg"
                      className="sidebar-inner"
                      alt="settings"
                    />
                    <div className="sidebar-head sidebar-inner">Settings</div>
                  </div>
                </div>
              </Link>
            </div>
          </div>
        </div>
      </Router>
    );
  }
}
