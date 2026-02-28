import React from "react";
import { render, screen } from "@testing-library/react";
import { vi, describe, it, expect } from 'vitest';
import { ThemeProvider, createTheme } from "@mui/material/styles";
import {EmptyShows} from ".";

// Mock the styles
vi.mock("./styles/emptyShowStyles", () => ({
    __esModule: true,
    default: () => ({
        emptyShowsLayout: {},
        emptyShowsIcon: {},
        emptyShowsContainer: {},
    }),
}));

// Create a mock theme for testing
const mockTheme = createTheme();

describe("Basic rendering", () => {
    it("Should render with message", () => {
        const testMessage = "Test Message";

        render(
            <ThemeProvider theme={mockTheme}>
                <EmptyShows emptyShowsMessage={testMessage}/>
            </ThemeProvider>
        );

        expect(screen.getByText(testMessage)).toBeInTheDocument();
    });
});
