import axios from "axios";
import {urls} from "../config/env-config";

const tokenKey = 'skyfox_token';

export const authHeader = () => {
    return {
        headers: {
            Authorization: 'Basic ' + localStorage.getItem(tokenKey)
        }
    };
}

export const login = async (username, password) => {
    const token = authBasic(username, password);
    const config = {
        headers: {
            Authorization: 'Basic ' + token
        }
    };
    const response = await axios.get(`${urls.service}/login`, config);
    const userDetails = response.data;
    localStorage.setItem(tokenKey, token)
    return userDetails;
}

export const isLoggedIn = () => {
    return localStorage.getItem(tokenKey) !== null;
}

export const logout = () => {
    localStorage.removeItem(tokenKey);
};

const authBasic = (username, password) => {
    return window.btoa(username + ':' + password);
}
