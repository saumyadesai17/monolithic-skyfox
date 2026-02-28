import React, { useState } from "react";
import { Typography } from "@mui/material";
import styles from "../styles/loginStyles"

export default (onLogin) => {
    const classes = styles();
    const [showError, setShowError] = useState(false);

    const errorMessage = () => {
        if (showError) {
            return (
                <Typography variant="body1" color="error" className={classes.loginErrorMessage}>
                    Login failed
                </Typography>
            )
        }
    };

    const handleLogin = async (values) => {
        const {username, password} = values;
        try {
            await onLogin(username, password);
            setShowError(false);
        } catch (err) {
            if (err.response && err.response.status === 401) {
                setShowError(true);
            } else {
                throw err;
            }
        }
    };

    return {
        errorMessage: errorMessage,
        handleLogin: handleLogin
    };
};
