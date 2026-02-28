import React from 'react'
import { Dialog, DialogContent, Typography, Alert } from "@mui/material";
import styles from "./styles/customerDetailsDialogStyles"

const BookingConfirmation = ({bookingConfirmation, showConfirmation}) => {
    const classes = styles();
    return (
      <Dialog open={showConfirmation}>
            <Alert severity="success">
                Seats booked successfully!
            </Alert>
            <Typography variant="h6" className={classes.dialogHeader}>
                Booking Confirmation
            </Typography>
            <DialogContent>
                <Typography variant="body1" display="block" gutterBottom>
                    Booking id : {bookingConfirmation.id}
                </Typography>
                <Typography variant="body1" display="block" gutterBottom>
                    Show Date: {bookingConfirmation.showDate}
                </Typography>
                <Typography variant="body1" display="block" gutterBottom>
                    Show start time: {bookingConfirmation.startTime}
                </Typography>
                <Typography variant="body1" display="block" gutterBottom>
                    Customer Name: {bookingConfirmation.customerName}
                </Typography>
                <Typography variant="body1" display="block" gutterBottom>
                    Amount Paid: {bookingConfirmation.amountPaid}
                </Typography>
                <Typography variant="body1" display="block" gutterBottom>
                    Number of seats booked: {bookingConfirmation.noOfSeats}
                </Typography>
            </DialogContent>
        </Dialog>
    )
}

export default BookingConfirmation;
