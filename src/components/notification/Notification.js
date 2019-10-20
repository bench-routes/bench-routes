/* eslint-disable jsx-a11y/no-noninteractive-element-interactions */
import React from 'react';

export default class Notification extends React.Component {
  componentDidMount() {
    if (global.showNotificationScreen) {
      document.getElementById('notification').style.width = '40%';
      document.getElementById('notification').style.display = 'block';
    }
  }

  updateShowNotificationsScreen = () => {
    global.showNotificationSection = !global.showNotificationSection;
    if (global.showNotificationSection) {
      document.getElementById('notification').style.width = '40%';
    } else {
      document.getElementById('notification').style.width = '0%';
    }
  };

  render() {
    return (
      <div>
        <div className="notification-icon">
          <img
            src="assets/icons/notify-icon.svg"
            alt="notification"
            onClick={() => this.updateShowNotificationsScreen()}
            onKeyDown={(e) => {
              if (e.keyCode === 13) {
                this.updateShowNotificationsScreen();
              }
            }}
          />
        </div>
        <div id="notification" className="notification-screen">
          <div
            style={{
              display: 'inline-flex',
              padding: '1% 0% 3% 2%',
              borderBottom: '1px solid #fff',
              width: '100%',
            }}
          >
            <img
              src="assets/icons/cross.svg"
              alt="collapse notifications"
              onClick={() => this.updateShowNotificationsScreen()}
              onKeyDown={(e) => {
                if (e.keyCode === 13) {
                  this.updateShowNotificationsScreen();
                }
              }}
            />
            <div>Notifications</div>
          </div>
          <div className="notification-messages">this is notification screen</div>
        </div>
      </div>
    );
  }
}
