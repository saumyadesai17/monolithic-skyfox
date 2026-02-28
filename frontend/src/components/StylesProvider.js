import React from 'react';
import { ThemeProvider, createTheme } from '@mui/material/styles';
import { StylesProvider as MuiStylesProvider } from '@mui/styles';
import Theme from '../Theme';

// Create a wrapper component that provides both ThemeProvider and StylesProvider
const StylesProvider = ({ children }) => {
  return (
    <MuiStylesProvider>
      <ThemeProvider theme={Theme}>
        {children}
      </ThemeProvider>
    </MuiStylesProvider>
  );
};

export default StylesProvider;
