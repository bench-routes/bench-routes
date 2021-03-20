import { blue, grey } from '@material-ui/core/colors';
import {
  createMuiTheme,
  responsiveFontSizes,
  ThemeProvider
} from '@material-ui/core/styles';
import React, { useState } from 'react';
import './App.css';
import './assets/bootstrap.min.css';
import BaseLayout from './layouts/BaseLayout';

function App() {
  const [darkMode, setDarkMode] = useState<boolean>(
    localStorage.getItem('dark-mode') === 'true' ? true : false
  );
  let theme = createMuiTheme({
    palette: {
      // Provides you with all
      // shades of whites
      type: darkMode ? 'dark' : 'light',
      primary: blue
      // secondary: <Color>,
    },
    typography: {
      fontFamily: [
        '-apple-system',
        'BlinkMacSystemFont',
        '"Segoe UI"',
        'Roboto',
        '"Helvetica Neue"',
        'Arial',
        'sans-serif',
        '"Apple Color Emoji"',
        '"Segoe UI Emoji"',
        '"Segoe UI Symbol"'
      ].join(','),
      fontSize: 12
    },
    overrides: {
      MuiAppBar: {
        colorPrimary: {
          backgroundColor: darkMode ? grey[700] : blue[500]
        }
      }
    }
  });
  theme = responsiveFontSizes(theme);
  const toggleDarkMode = () => {
    setDarkMode(!darkMode);
  };
  return (
    <ThemeProvider theme={theme}>
      <BaseLayout toggleDarkMode={toggleDarkMode} darkMode={darkMode} />
    </ThemeProvider>
  );
}

export default App;
