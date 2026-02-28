import { Navigate, useLocation } from "react-router-dom";
import React from "react";
import PropTypes from "prop-types";

const ProtectedRoute = ({ children, isAuthenticated }) => {
    const location = useLocation();
    
    if (!isAuthenticated) {
        // Redirect to login page but save the current location
        return <Navigate to="/login" state={{ from: location }} replace />;
    }
    
    return children;
};

ProtectedRoute.propTypes = {
    children: PropTypes.node.isRequired,
    isAuthenticated: PropTypes.bool.isRequired
};

export default ProtectedRoute;
