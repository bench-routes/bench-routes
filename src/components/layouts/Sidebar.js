import React from 'react';
import './layouts.style.css';

export default class Sidebar extends React.Component {
    render() {
        return (
            <div className='sidebar'>
                <div>
                <div className='sidebar-inner'>
                    <img src='assets/icons/dashboard-icon.svg' className='sidebar-inner' alt='dashboard' />
                    <div className='sidebar-head sidebar-inner'>Dashboard</div>
                </div>
                </div>
                <div>
                <div className='sidebar-inner'>
                    <img src='assets/icons/monitoring-icon.svg' className='sidebar-inner' alt='monitoring' />
                    <div className='sidebar-head sidebar-inner'>Monitoring</div>
                </div>
                </div>
                <div>
                <div className='sidebar-inner'>
                    <img src='assets/icons/bench-icon.svg' className='sidebar-inner' alt='Benchmarks' />
                    <div className='sidebar-head sidebar-inner'>Benchmarks</div>
                </div>
                </div>
            </div>
        );
    }
}