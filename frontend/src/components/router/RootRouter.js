import { BrowserRouter as Router, Routes, Route, Navigate } from "react-router-dom";
import React from "react";
import Shows from "../shows/Shows";
import ErrorOutlineIcon from '@mui/icons-material/ErrorOutline';
import BlockIcon from '@mui/icons-material/Block';
import { Error } from "../common";
import { Login, ProtectedRoute } from "../login";
import PropTypes from "prop-types";
import moment from "moment";

const RootRouter = ({ isAuthenticated, onLogin }) => {
    const todayDate = moment().format("YYYY-MM-DD");

    return (
        <Router>
            <Routes>
                <Route path="/" element={<Navigate to={`/shows?date=${todayDate}`} />} />
                
                <Route 
                    path="/shows" 
                    element={
                        <ProtectedRoute isAuthenticated={isAuthenticated}>
                            <Shows />
                        </ProtectedRoute>
                    } 
                />

                <Route 
                    path="/login" 
                    element={<Login isAuthenticated={isAuthenticated} onLogin={onLogin} />} 
                />

                <Route 
                    path="/error" 
                    element={<Error errorIcon={ErrorOutlineIcon} errorMessage={"Oops..Something went wrong"} />} 
                />

                <Route 
                    path="*" 
                    element={<Error errorIcon={BlockIcon} errorMessage={"Not Found"} />} 
                />
            </Routes>
        </Router>
    );
};

RootRouter.propTypes = {
    isAuthenticated: PropTypes.bool.isRequired,
    onLogin: PropTypes.func.isRequired
};

export default RootRouter;
