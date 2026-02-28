import React from "react";
import {fireEvent, render} from "@testing-library/react";
import SeatSelectionDialog from "./SeatSelectionDialog";
import { vi, describe, it, expect } from 'vitest';
import { ThemeProvider, createTheme } from "@mui/material/styles";

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
vi.mock("./styles/seatSelectionDialogStyles", () => ({
    __esModule: true,
    default: () => ({
        dialogTitle: {},
        dialogContent: {},
        dialogActions: {},
        formControl: {},
        selectEmpty: {},
        submitButton: {},
        movieDetails: {},
    }),
}));

vi.mock("./CustomerDetailsDialog", () => ({
    __esModule: true,
    default: vi.fn(({open}) => <div>Customer Details is {open ? "open" : "closed"}</div>)
}));

describe("Basic rendering and functionality", () => {
    const openDialog = true;
    const onClose = vi.fn();
    const updateShowRevenue = vi.fn();

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
        const {queryByText, queryAllByText, queryByDisplayValue} = render(
            <ThemeProvider theme={mockTheme}>
                <SeatSelectionDialog 
                    selectedShow={selectedShow}
                    open={openDialog} 
                    onClose={onClose}
                    updateShowsRevenue={updateShowRevenue}
                />
            </ThemeProvider>
        );

        expect(queryByText(selectedShow.movie.name)).toBeTruthy();
        expect(queryByText(selectedShow.movie.plot)).toBeTruthy();
        expect(queryByText(selectedShow.movie.duration)).toBeTruthy();
        expect(queryAllByText("Seats").length).toBeGreaterThan(0);
        expect(queryByDisplayValue("1")).toBeTruthy();
    });

    it("Should display total cost when number of seats is selected", () => {
        const {queryByText, getByDisplayValue} = render(
            <ThemeProvider theme={mockTheme}>
                <SeatSelectionDialog 
                    selectedShow={selectedShow}
                    open={openDialog} 
                    onClose={onClose}
                    updateShowsRevenue={updateShowRevenue}
                />
            </ThemeProvider>
        );

        expect(queryByText("₹150.00")).toBeTruthy();
        fireEvent.change(getByDisplayValue("1"), {target: {value: '2'}});

        expect(queryByText("₹300.00")).toBeTruthy();
    });

    it("Should display customer details input on next", () => {
        const {getByText} = render(
            <ThemeProvider theme={mockTheme}>
                <SeatSelectionDialog 
                    selectedShow={selectedShow} 
                    open={openDialog}
                    onClose={onClose}
                    updateShowsRevenue={updateShowRevenue}
                />
            </ThemeProvider>
        );

        expect(getByText("Customer Details is closed")).toBeTruthy();

        fireEvent.click(getByText("Next"));

        expect(getByText("Customer Details is open")).toBeTruthy();
    });
});
