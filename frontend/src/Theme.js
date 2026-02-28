import { green, red } from '@mui/material/colors';
import { createTheme } from '@mui/material/styles';

// A custom theme for this app
const theme = createTheme({
    palette: {
        primary: {
            main: '#556cd6',
            contrastText: '#fff'
        },
        secondary: {
            main: '#19857b',
        },
        success: {
            main: green['800']
        },
        error: {
            main: red.A400,
        },
        background: {
            default: '#fff',
        },
    },
});

export default theme;
