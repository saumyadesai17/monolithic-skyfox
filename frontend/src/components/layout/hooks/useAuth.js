import {useEffect, useState} from "react";
import {isLoggedIn, login, logout} from "../../../helpers/authService";

export default () => {
    const [isAuthenticated, setIsAuthenticated] = useState(false);
    useEffect(() => {
        setIsAuthenticated(isLoggedIn());
    }, []);

    const handleLogin = async (username, password) => {
        const userDetails = await login(username, password);
        setIsAuthenticated(true);
        return userDetails;
    };

    const handleLogout = () => {
        logout();
        setIsAuthenticated(false);
    };

    return {
        isAuthenticated: isAuthenticated,
        handleLogin: handleLogin,
        handleLogout: handleLogout
    };
}
