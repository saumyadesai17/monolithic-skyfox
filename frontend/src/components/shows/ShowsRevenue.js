import React from "react";
import { CircularProgress, Typography } from "@mui/material";
import styles from "./styles/showsRevenueStyles"
import PropTypes from "prop-types";
import {INR_SYMBOL} from "../../Constants";

const ShowsRevenue = ({showsRevenue, showsRevenueLoading}) => {
    const classes = styles();

    return (
        <>
            {
                showsRevenueLoading
                    ? (
                        <div className={classes.showsRevenueLoadingSpinner}>
                            <CircularProgress color="primary"/>
                        </div>
                    )
                    : (
                        <Typography variant="h5" color="secondary" className={classes.showsRevenueContainer}>
                            Revenue: {`${INR_SYMBOL}${showsRevenue}`}
                        </Typography>
                    )
            }
        </>
    );
};

ShowsRevenue.propTypes = {
    showsRevenue: PropTypes.number.isRequired,
    showsRevenueLoading: PropTypes.bool.isRequired
};

export default ShowsRevenue;
