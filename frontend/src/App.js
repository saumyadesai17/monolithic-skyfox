import React from 'react';
import { CssBaseline, ThemeProvider } from '@mui/material';
import Layout from "./components/layout/Layout";
import Theme from './Theme';

const App = () => {
    return (
        <ThemeProvider theme={Theme}>
            <CssBaseline />
            <Layout />
        </ThemeProvider>
    );
};

export default App;
