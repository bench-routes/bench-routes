import React from 'react';

export default class Notification extends React.Component {

  updateShowNotificationsScreen = () => {
    console.log('updating');
    global.showNotificationSection = !global.showNotificationSection;
    console.log(global.showNotificationSection)
    if (global.showNotificationSection) {
      document.getElementById('notification').style.width = '40%';
    } else {
      document.getElementById('notification').style.width = '0%';
    }
  }

  componentDidMount() {
    console.log('called')
    if (global.showNotificationScreen) {
      document.getElementById('notification').style.width = '40%';
      document.getElementById('notification').style.display = 'block';
    }
  }

  render() {
    return (
      <div>
        <div className='notification-icon'>
          <img src='assets/icons/notify-icon.svg' alt='notification' onClick={() => this.updateShowNotificationsScreen()} />
        </div>
        <div id='notification' className='notification-screen'>
          <div style={{ display: 'inline-flex', padding: `1% 0% 3% 2%`, borderBottom: '1px solid #fff', width: '100%' }}>
            <img src='assets/icons/cross.svg' alt='collapse notifications' onClick={() => this.updateShowNotificationsScreen()} />
            <div>
              Notifications
            </div>
          </div>
          <div className='notification-messages'>
            this is notification screen
          </div>
        </div>
      </div>
    );
  }
}