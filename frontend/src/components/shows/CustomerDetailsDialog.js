import { Button, Dialog, DialogContent, Typography, Snackbar, Alert } from "@mui/material";
import React, { useState } from "react";
import styles from "./styles/customerDetailsDialogStyles"
import bookingService from "./services/bookingService";
import PropTypes from "prop-types";
import moment from "moment";
import {object, string} from "yup";
import {Form, Formik} from "formik";
import FormikTextField from "../formik/FormikTextField";
import BookingConfirmation from "./BookingConfirmation";

const CustomerDetailsDialog = ({seats, selectedShow, updateShowsRevenue, open, onClose}) => {
    const [success, setSuccess] = useState(null);
    const [showConfirmation, setShowConfirmation] = useState(false);
    const [bookingConfirmation, setBookingConfirmation] = useState({});
    const classes = styles();

    const initialValues = {
        name: "",
        phoneNumber: ""
    };

    const formSchema = object({
        name: string("Enter name")
            .required("Name is required"),
        phoneNumber: string("Enter phone number")
            .required("Phone number is required")
            .matches(/^\d{10}$/, "Phone number should be 10 digits")
    });

    const bookShow = async (values) => {
        const today = moment().format("YYYY-MM-DD");
        const payload = {
            date: today,
            showId: selectedShow.id,
            customer: {
                name: values.name,
                phoneNumber: values.phoneNumber
            },
            noOfSeats: Number(seats)
        };

        try {
            const response = await bookingService.create(payload);
            setSuccess(true);
            updateShowsRevenue();
            setBookingConfirmation(response.data)
            setShowConfirmation(true);
        } catch {
            setSuccess(false);
        } finally {
            onClose();
        }
    };
    return (
      <>
            <Dialog open={open} onClose={onClose} maxWidth={false}>
                <Typography variant="h6" className={classes.dialogHeader}>
                    Enter Customer Details
                </Typography>
                <Formik validationSchema={formSchema} initialValues={initialValues} onSubmit={bookShow}>
                    {
                        ({isValid}) => {
                            return (
                                <Form>
                                    <DialogContent className={classes.dialogContent}>
                                        <FormikTextField
                                            required
                                            margin="dense"
                                            inputProps={{"data-testid": "name"}}
                                            name="name"
                                            label="Name"
                                            fullWidth
                                            autoComplete='off'
                                            autoFocus
                                        />
                                        <FormikTextField
                                            required
                                            margin="dense"
                                            inputProps={{"data-testid": "phoneNumber"}}
                                            name="phoneNumber"
                                            label="Phone Number"
                                            fullWidth
                                            autoComplete='off'
                                        />
                                        <Button type="submit" disabled={!isValid} color="primary" variant="contained"
                                                className={classes.bookShowButton} data-testid="bookButton">
                                            Book
                                        </Button>
                                    </DialogContent>
                                </Form>
                            );
                        }
                    }
                </Formik>
            </Dialog>

            <BookingConfirmation bookingConfirmation={bookingConfirmation} showConfirmation={showConfirmation} />

            <Snackbar open={success === false} autoHideDuration={2000} onClose={() => setSuccess(null)}>
                <Alert severity="error">
                    Sorry, seats could not be booked!
                </Alert>
            </Snackbar>
        </>
    )
}

CustomerDetailsDialog.propTypes = {
    selectedShow: PropTypes.object.isRequired,
    seats: PropTypes.string.isRequired,
    updateShowsRevenue: PropTypes.func.isRequired,
    open: PropTypes.bool.isRequired,
    onClose: PropTypes.func.isRequired
};

export default CustomerDetailsDialog;
