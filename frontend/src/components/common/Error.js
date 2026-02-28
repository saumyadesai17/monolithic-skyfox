import React from "react";
import { Typography } from "@mui/material";
import styles from "./styles/errorStyles"
import PropTypes from "prop-types";

const Error = ({errorIcon: ErrorIcon, errorMessage}) => {
    const classes = styles();

    return (
        <div className={classes.errorContent}>
            <ErrorIcon className={classes.errorIcon}/>
            <Typography variant='h4'>
                {errorMessage}
            </Typography>
        </div>
    );
};

Error.propTypes = {
    errorIcon: PropTypes.elementType.isRequired,
    errorMessage: PropTypes.string.isRequired,
};

export default Error;
