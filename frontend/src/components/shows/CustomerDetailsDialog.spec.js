import React from "react";
import {fireEvent, render, waitFor} from "@testing-library/react";
import CustomerDetailsDialog from "./CustomerDetailsDialog";
import bookingService from "./services/bookingService";
import moment from "moment";
import { ThemeProvider, createTheme } from "@mui/material/styles";

// Mock the BookingConfirmation component
vi.mock("./BookingConfirmation", () => ({
    __esModule: true,
    default: vi.fn(({ bookingConfirmation, onClose }) => (
        <div data-testid="booking-confirmation">
            Booking Confirmation
        </div>
    ))
}));

// Create a mock theme
const mockTheme = createTheme({
    palette: {
        primary: {
            main: '#1976d2',
        },
        secondary: {
            main: '#dc004e',
        },
        background: {
            paper: '#fff',
        },
    },
    spacing: (factor) => `${0.25 * factor}rem`,
    shadows: Array(25).fill('none'),
});

// Mock styles
vi.mock("./styles/customerDetailsDialogStyles", () => ({
    __esModule: true,
    default: () => ({
        dialogTitle: {},
        dialogContent: {},
        dialogActions: {},
        textField: {},
        submitButton: {},
    }),
}));

vi.mock("./services/bookingService", () => ({
    __esModule: true,
    default: {
        create: vi.fn()
    }
}));

vi.mock("moment");

describe("Basic rendering and functionality", () => {
    const open = true;
    const onClose = vi.fn();
    const updateShowsRevenue = vi.fn();
    const selectedShow = {
        id: 1,
        cost: 150,
        movie: {
            name: "Movie 1",
            plot: "Suspense movie",
            duration: "1hr 30m"
        },
        slot: {startTime: "start time 1"}
    };

    it("Should display the show info", () => {
        const {queryByText} = render(
            <ThemeProvider theme={mockTheme}>
                <CustomerDetailsDialog 
                    seats={"2"} 
                    selectedShow={selectedShow} 
                    open={open}
                    onClose={onClose}
                    updateShowsRevenue={updateShowsRevenue}
                />
            </ThemeProvider>
        );

        expect(queryByText("Enter Customer Details")).toBeTruthy();
        expect(queryByText("Name")).toBeTruthy();
        expect(queryByText("Phone Number")).toBeTruthy();
    });

    it("Should call booking service api to create booking on submit", async () => {
        const {getByTestId} = render(
            <ThemeProvider theme={mockTheme}>
                <CustomerDetailsDialog 
                    seats={"2"} 
                    selectedShow={selectedShow}
                    open={open}
                    onClose={onClose}
                    updateShowsRevenue={updateShowsRevenue}
                />
            </ThemeProvider>
        );

        // Mock moment to return a date formatter function
        const testFormat = vi.fn();
        testFormat.mockImplementation((format) => {
            if (format === "YYYY-MM-DD") {
                return "2020-06-19";
            }
            return null;
        });

        moment.mockReturnValue({
            format: testFormat
        });

        // Mock bookingService.create to resolve with an empty string
        bookingService.create.mockImplementation((payload) => {
            if (payload && 
                payload.customer && 
                payload.customer.name === "Name" && 
                payload.customer.phoneNumber === "1234567890") {
                return Promise.resolve("");
            }
            return Promise.reject(new Error("Unexpected payload"));
        });

        fireEvent.change(getByTestId("name"), {
            target: {
                value: "Name"
            }
        });

        fireEvent.change(getByTestId("phoneNumber"), {
            target: {
                value: "1234567890"
            }
        });

        fireEvent.click(getByTestId("bookButton"));

        const expectedPayload = {
            "customer": {"name": "Name", "phoneNumber": "1234567890"},
            "date": "2020-06-19",
            "noOfSeats": 2,
            "showId": 1
        };

        await waitFor(() => {
            expect(bookingService.create).toHaveBeenCalledTimes(1);
            expect(bookingService.create).toHaveBeenCalledWith(expectedPayload);
        });
    });
});
