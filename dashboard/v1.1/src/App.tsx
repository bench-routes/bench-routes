import React from 'react';
import './App.css';
import {
  createMuiTheme,
  ThemeProvider,
  responsiveFontSizes
} from '@material-ui/core/styles';
import { grey } from '@material-ui/core/colors';

let theme = createMuiTheme({
  palette: {
    // Provides you with all
    // shades of whites
    primary: grey
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
      <div className="App">Bench-routes ui</div>
    </ThemeProvider>
  );
}

export default App;
