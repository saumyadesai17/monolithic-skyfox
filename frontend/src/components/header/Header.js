import React from "react";
import { AppBar, Toolbar, Typography, Box } from "@mui/material";
import MovieIcon from '@mui/icons-material/Movie';
import ExitToAppIcon from '@mui/icons-material/ExitToApp';
import styles from "./styles/headerStyles";
import PropTypes from "prop-types";

const Header = ({onLogout, isAuthenticated}) => {
    const logoutSection = () => {
        if (isAuthenticated) {
            return (
                <Box sx={styles.logoutLink} onClick={onLogout}>
                    <ExitToAppIcon/>
                    <Typography sx={styles.headerLogo} variant="body1">
                        Logout
                    </Typography>
                </Box>
            );
        }
    };

    return (
        <AppBar position={"sticky"}>
            <Toolbar sx={styles.toolbar}>
                <Box component="a" href="/" sx={styles.headerLink}>
                    <MovieIcon sx={styles.cinemaLogoIcon}/>
                    <Typography sx={styles.headerLogo} variant="h5">
                        SkyFox Cinema
                    </Typography>
                </Box>
                {logoutSection()}
            </Toolbar>
        </AppBar>
    );
};

Header.propTypes = {
    onLogout: PropTypes.func.isRequired,
    isAuthenticated: PropTypes.bool.isRequired
};

export default Header;
