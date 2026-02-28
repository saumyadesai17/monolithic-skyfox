import React from "react";
import { render, screen } from "@testing-library/react";
import { vi, describe, it, expect } from 'vitest';
import { ThemeProvider, createTheme } from "@mui/material/styles";
import { red } from "@mui/material/colors";
import Error from "./Error";

// Mock the styles
vi.mock("./styles/errorStyles", () => ({
    __esModule: true,
    default: () => ({
        errorContent: {},
        errorIcon: {},
    }),
}));

// Create a mock theme for testing
const mockTheme = createTheme({
  palette: {
    error: {
      main: red.A400,
    },
  },
});

describe("Basic rendering", () => {
    it("Should render with icon and message", () => {
        const testErrorMessage = "Test Error";
        const TestErrorIcon = () => <span data-testid="test-error-icon" />;

        render(
            <ThemeProvider theme={mockTheme}>
                <Error errorIcon={TestErrorIcon} errorMessage={testErrorMessage}/>
            </ThemeProvider>
        );
        
        expect(screen.getByTestId("test-error-icon")).toBeInTheDocument();
        expect(screen.getByText(testErrorMessage)).toBeInTheDocument();
    });
});
