import React from "react";
import { render, screen } from "@testing-library/react";
import ShowsRevenue from "./ShowsRevenue";

// Mock the MUI components
vi.mock("@mui/material", () => ({
    CircularProgress: ({ color }) => <div data-testid="circular-progress" data-color={color}>Loading...</div>,
    Typography: ({ children, variant, color, className }) => (
        <div data-testid="typography" data-variant={variant} data-color={color} className={className}>
            {children}
        </div>
    )
}));

describe("Basic rendering", () => {
    it("Should show revenue if not loading", () => {
        render(<ShowsRevenue showsRevenue={549.99} showsRevenueLoading={false} />);

        // Check that the revenue is displayed
        expect(screen.getByText("Revenue: ₹549.99")).toBeInTheDocument();
    });

    it("Should display spinner if loading", () => {
        render(<ShowsRevenue showsRevenue={0} showsRevenueLoading={true} />);

        // Check that the circular progress is displayed
        expect(screen.getByTestId("circular-progress")).toBeInTheDocument();
        expect(screen.getByTestId("circular-progress")).toHaveAttribute("data-color", "primary");
    });
});
