import React from 'react';
import Header from './Header';
import Sidebar from './Sidebar';

export default function BaseLayout(props) {
  return (
    <div className="App">
      <div className="sidebar">
        <Sidebar />
      </div>
      <div className="inner-component">
        <Header />
        {props.children}
      </div>
    </div>
  );
}
