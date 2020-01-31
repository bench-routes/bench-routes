import React, { FC, useState } from 'react';
import { HashRouter as Router, Link } from 'react-router-dom';
import './layouts.style.css';
import { Collapse } from 'reactstrap';

const Sidebar: FC<{}> = () => {
  const [showSubmenu, setShowSubmenu] = useState(false);
  const toggleBenchmarkSubmenu = () => {
    setShowSubmenu(!showSubmenu);
  };

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
            <div className="sidebar-inner">
              <img
                src="assets/icons/monitoring-icon.svg"
                className="sidebar-inner"
                alt="monitoring"
              />
              <div className="sidebar-head sidebar-inner">Monitoring</div>
            </div>
          </Link>

          <div onClick={() => toggleBenchmarkSubmenu()}>
            <div className="sidebar-inner benchmarking">
              <img
                src="assets/icons/bench-icon.svg"
                className="sidebar-inner"
                alt="Tests"
              />
              <div className="sidebar-head sidebar-inner">Tests</div>
            </div>
          </div>

          <div className="benchmark-submenu">
            <Collapse isOpen={showSubmenu}>
              <div key="compulsory_transition_key">
                <Link to="/ping" style={{ textDecoration: 'none' }}>
                  <div>
                    <div className="sidebar-inner">
                      <img
                        src="assets/icons/ping-meter.svg"
                        className="sidebar-submenu-inner"
                        alt="Benchmarks"
                      />
                      <div className="sidebar-head sidebar-inner">Ping</div>
                    </div>
                  </div>
                </Link>
                <Link to="/floodping" style={{ textDecoration: 'none' }}>
                  <div>
                    <div className="sidebar-inner">
                      <img
                        src="assets/icons/flood-icon.png"
                        className="sidebar-submenu-inner"
                        alt="Flood-Ping"
                      />
                      <div className="sidebar-head sidebar-inner">
                        Floodping
                      </div>
                    </div>
                  </div>
                </Link>
                <Link to="/jitter" style={{ textDecoration: 'none' }}>
                  <div>
                    <div className="sidebar-inner">
                      <img
                        src="assets/icons/jitter-icon.png"
                        className="sidebar-submenu-inner"
                        alt="Jitter"
                      />
                      <div className="sidebar-head sidebar-inner">Jitter</div>
                    </div>
                  </div>
                </Link>
              </div>
            </Collapse>
          </div>

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
};

export default Sidebar;
