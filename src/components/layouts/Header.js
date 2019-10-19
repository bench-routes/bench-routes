import React from 'react';

export default class Header extends React.Component {
  render() {
    return (
      <header>
        <div className='logo-name'>
          Bench-routes
        </div>
        
        <div className='notification-icon'>
          <img src='assets/icons/notify-icon.svg' alt='notification' />
        </div>
      </header>
    );
  }
}