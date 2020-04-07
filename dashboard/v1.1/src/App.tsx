import { blue } from '@material-ui/core/colors';
import {
  createMuiTheme,
  responsiveFontSizes,
  ThemeProvider
} from '@material-ui/core/styles';
import React from 'react';
import './App.css';
import BaseLayout from './layouts/BaseLayout';
import Navigator from './router/Navigation';

let theme = createMuiTheme({
  palette: {
    // Provides you with all
    // shades of whites
    primary: blue
    // secondary: <Color>,
  },
  typography: {
    fontFamily: ['Lato', 'Raleway'].join(','),
    fontSize: 12
  }
});
theme = responsiveFontSizes(theme);

function App() {
  return (
    <ThemeProvider theme={theme}>
      <BaseLayout>
        <Navigator />
      </BaseLayout>
    </ThemeProvider>
  );
}

export default App;
