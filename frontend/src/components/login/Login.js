import React, { useEffect } from "react";
import { useLocation, useNavigate } from "react-router-dom";
import { Form, Formik } from "formik";
import { FormikTextField } from "../formik";
import { Button } from "@mui/material";
import styles from "./styles/loginStyles"
import PropTypes from "prop-types";
import useLogin from "./hooks/useLogin";
import { formSchema, initialValues } from "./services/loginFormService";

const Login = ({ isAuthenticated, onLogin }) => {
    const classes = styles();
    const location = useLocation();
    const navigate = useNavigate();
    
    // In react-router v6, the state structure is different
    const from = location.state?.from?.pathname || "/";
    const { errorMessage, handleLogin } = useLogin(onLogin);

    useEffect(() => {
        if (isAuthenticated) {
            navigate(from, { replace: true });
        }
    }, [isAuthenticated, navigate, from]);

    return (
        <div className={classes.loginContainer}>
            <Formik initialValues={initialValues}
                    onSubmit={handleLogin}
                    validationSchema={formSchema}>
                {
                    (props) => {
                        const {
                            isValid,
                        } = props;
                        return (
                            <Form className={classes.loginForm}>
                                <FormikTextField
                                    required
                                    margin="dense"
                                    name="username"
                                    label="Username"
                                />
                                <FormikTextField
                                    required
                                    type="password"
                                    margin="dense"
                                    name="password"
                                    label="Password"
                                />
                                {
                                    errorMessage()
                                }
                                <Button
                                    variant="contained"
                                    type="submit"
                                    disabled={!isValid}
                                    color="primary"
                                    className={classes.loginButton}
                                >
                                    Login
                                </Button>
                            </Form>
                        );
                    }
                }
            </Formik>
        </div>
    );
}

Login.propTypes = {
    isAuthenticated: PropTypes.bool.isRequired,
    onLogin: PropTypes.func.isRequired
};

export default Login;
