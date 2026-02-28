import {object, string} from "yup";

export const initialValues = {
    username: '',
    password: ''
};

export const formSchema = object({
    username: string("Enter username")
        .required("Username is required"),
    password: string("Enter password")
        .required("Password is required")
});
