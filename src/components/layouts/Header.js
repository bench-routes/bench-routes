/* eslint-disable jsx-a11y/no-noninteractive-element-interactions */
import React, { useState } from 'react';
import Notification from '../notification/Notification';

export default function Header() {
  const [showNotification, setShowNotification] = useState(false);

  const updateShowNotificationsScreen = () => {
    setShowNotification(!showNotification);
  };

  return (
    <>
      <Notification
        showNotification={showNotification}
        updateShowNotificationsScreen={updateShowNotificationsScreen}
      />
      <header>
        <div className="logo-name">Bench-routes</div>
        <div className="notification-icon">
          <img
            src="assets/icons/notify-icon.svg"
            alt="notification"
            onClick={() => updateShowNotificationsScreen()}
            onKeyDown={(e) => {
              if (e.keyCode === 13) {
                updateShowNotificationsScreen();
              }
            }}
          />
        </div>
      </header>
    </>
  );
}
