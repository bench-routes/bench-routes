import React from 'react';
import './App.css';
import Header from './components/layouts/Header';
import Sidebar from './components/layouts/Sidebar';

function App() {
  return (
    <div className="App">
      <Sidebar className="sidebar" />
      <div className="main">
        <Header />
      </div>
    </div>
  );
}

export default App;
