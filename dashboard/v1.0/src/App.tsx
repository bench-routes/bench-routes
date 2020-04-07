import React from 'react';
import './App.css';
import BaseLayout from './components/layouts/BaseLayout';
import Navigator from './router/Navigation';

function App() {
  return (
    <BaseLayout>
      <Navigator />
    </BaseLayout>
  );
}

export default App;
