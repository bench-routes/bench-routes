import { render } from '@testing-library/react';
import React from 'react';
import App from './App';

test('renders without crashing', () => {
  const { findByTestId } = render(<App />);
  const element = findByTestId('App');
  expect(element).toBeTruthy();
});
