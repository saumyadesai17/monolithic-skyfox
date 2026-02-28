import React from "react";
import {fireEvent, render, screen} from "@testing-library/react";
import Shows from "./Shows";
import {dateFromSearchString, nextDateLocation, previousDateLocation} from "./services/dateService";
import useShows from "./hooks/useShows";
import SeatSelectionDialog from "./SeatSelectionDialog";
import useShowsRevenue from "./hooks/useShowsRevenue";
import ShowsRevenue from "./ShowsRevenue";
import * as reactRouterDom from "react-router-dom";
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
    zIndex: {
        drawer: 1200,
    },
    spacing: (factor) => `${0.25 * factor}rem`,
    shadows: Array(25).fill('none'),
});

// Mock styles
vi.mock("./styles/showsStyles", () => ({
    __esModule: true,
    default: () => ({
        cardHeader: {},
        showContainer: {},
        localMoviesIcon: {},
        showsHeader: {},
        backdrop: {},
        listRoot: {},
        price: {},
        slotTime: {},
        buttons: {},
        navigationButton: {},
        paper: {},
    }),
}));

// Mock react-router-dom hooks
vi.mock("react-router-dom", () => ({
    useLocation: vi.fn(),
    useNavigate: vi.fn()
}));

// Mock the ShowsRevenue component
vi.mock("./ShowsRevenue", () => ({
    __esModule: true,
    default: vi.fn(({ showsRevenue, showsRevenueLoading }) => (
        <div data-testid="shows-revenue" data-revenue={showsRevenue} data-loading={showsRevenueLoading}>
            Shows Revenue Component
        </div>
    ))
}));

vi.mock("./services/dateService", () => ({
    dateFromSearchString: vi.fn(),
    nextDateLocation: vi.fn(),
    previousDateLocation: vi.fn()
}));

vi.mock("./hooks/useShows", () => ({
    __esModule: true,
    default: vi.fn()
}));

vi.mock("./hooks/useShowsRevenue", () => ({
    __esModule: true,
    default: vi.fn()
}));

vi.mock("./SeatSelectionDialog", () => ({
    __esModule: true,
    default: vi.fn(() => <div>SeatSelection</div>)
}));

describe("Basic rendering and functionality", () => {
    let testNavigate;
    let testLocation;
    let testShowDate;

    beforeEach(() => {
        // Reset mocks
        vi.clearAllMocks();
        
        testNavigate = vi.fn();
        reactRouterDom.useNavigate.mockReturnValue(testNavigate);

        testLocation = {
            search: "testSearch"
        };
        reactRouterDom.useLocation.mockReturnValue(testLocation);

        testShowDate = {
            format: vi.fn()
        };

        // Mock dateFromSearchString to return testShowDate when called with "testSearch"
        dateFromSearchString.mockImplementation((search) => {
            if (search === "testSearch") {
                return testShowDate;
            }
            return null;
        });

        // Mock nextDateLocation to return "Next Location" when called with testLocation and testShowDate
        nextDateLocation.mockImplementation((location, date) => {
            if (location === testLocation && date === testShowDate) {
                return "Next Location";
            }
            return null;
        });

        // Mock previousDateLocation to return "Previous Location" when called with testLocation and testShowDate
        previousDateLocation.mockImplementation((location, date) => {
            if (location === testLocation && date === testShowDate) {
                return "Previous Location";
            }
            return null;
        });

        // Mock testShowDate.format to return "Show Date" when called with "Do MMM YYYY"
        testShowDate.format.mockImplementation((format) => {
            if (format === "Do MMM YYYY") {
                return "Show Date";
            }
            return null;
        });

        // Mock useShows to return shows data when called with testShowDate
        useShows.mockImplementation((date) => {
            if (date === testShowDate) {
                return {
                    showsLoading: false,
                    shows: [
                        {
                            id: 1,
                            cost: 150,
                            movie: {name: "Movie 1"},
                            slot: {startTime: "start time 1"}
                        }, {
                            id: 2,
                            cost: 160,
                            movie: {name: "Movie 2"},
                            slot: {startTime: "start time 2"}
                        }
                    ]
                };
            }
            return { showsLoading: true, shows: [] };
        });

        // Mock useShowsRevenue to return revenue data when called with testShowDate
        useShowsRevenue.mockImplementation((date) => {
            if (date === testShowDate) {
                return {
                    showsRevenue: 549.99,
                    showsRevenueLoading: false
                };
            }
            return { showsRevenue: 0, showsRevenueLoading: true };
        });
    });

    it("Should display the show info", () => {
        render(
            <ThemeProvider theme={mockTheme}>
                <Shows />
            </ThemeProvider>
        );

        expect(screen.getByText("Shows (Show Date)")).toBeInTheDocument();

        expect(screen.getByText("Movie 1")).toBeInTheDocument();
        expect(screen.getByText("start time 1")).toBeInTheDocument();
        expect(screen.getByText("₹150")).toBeInTheDocument();

        expect(screen.getByText("Movie 2")).toBeInTheDocument();
        expect(screen.getByText("start time 2")).toBeInTheDocument();
        expect(screen.getByText("₹160")).toBeInTheDocument();
    });

    it("Should push to history if next or previous clicked", () => {
        render(
            <ThemeProvider theme={mockTheme}>
                <Shows />
            </ThemeProvider>
        );

        const previousDayButton = screen.getByText("Previous Day");
        const nextDayButton = screen.getByText("Next Day");

        fireEvent.click(previousDayButton);
        fireEvent.click(nextDayButton);

        expect(testNavigate).toBeCalledTimes(2);
        expect(testNavigate).toHaveBeenNthCalledWith(1, "Previous Location");
        expect(testNavigate).toHaveBeenNthCalledWith(2, "Next Location");
    });

    it("Should display seat selection when a show is selected", () => {
        render(
            <ThemeProvider theme={mockTheme}>
                <Shows />
            </ThemeProvider>
        );

        expect(screen.queryByText("SeatSelectionDialog")).toBeNull();

        fireEvent.click(screen.getByText("Movie 1"));

        expect(screen.getByText("SeatSelection")).toBeInTheDocument();
    });

    it("Should display revenue when rendered", () => {
        render(
            <ThemeProvider theme={mockTheme}>
                <Shows />
            </ThemeProvider>
        );

        const showsRevenueElement = screen.getByTestId("shows-revenue");
        
        expect(showsRevenueElement).toBeInTheDocument();
        expect(showsRevenueElement.getAttribute("data-revenue")).toBe("549.99");
        expect(showsRevenueElement.getAttribute("data-loading")).toBe("false");
    });
});
