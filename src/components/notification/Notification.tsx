/* eslint-disable jsx-a11y/no-noninteractive-element-interactions */
import React from "react";

export default class Notification extends React.Component<{
  showNotification: boolean;
  updateShowNotificationsScreen: () => void;
}> {
  render() {
    return (
      <div
        className={`notification ${
          this.props.showNotification ? "display-notification" : "close-notification"
        }`}
      >
        <div className="notification-content">
          <div className="notification-header">
            <div>Notifications</div>
            <img
              src="assets/icons/cross.svg"
              alt="collapse notifications"
              onClick={() => this.props.updateShowNotificationsScreen()}
              onKeyDown={e => {
                if (e.keyCode === 13) {
                  return this.props.updateShowNotificationsScreen();
                }
              }}
            />
          </div>
        </div>
        <div className="notification-body">this is notification screen</div>
      </div>
    );
  }
}
