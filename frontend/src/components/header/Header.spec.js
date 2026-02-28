import React from "react";
import { render, screen, fireEvent } from "@testing-library/react";
import Header from "./Header";

// Mock the MUI components
vi.mock("@mui/material", () => ({
    AppBar: ({ children, position }) => <div data-testid="app-bar" data-position={position}>{children}</div>,
    Toolbar: ({ children, sx }) => <div data-testid="toolbar">{children}</div>,
    Typography: ({ children, variant, sx }) => <div data-testid="typography" data-variant={variant}>{children}</div>,
    Box: ({ children, sx, component, href, onClick }) => (
        <div 
            data-testid="box" 
            data-component={component} 
            data-href={href} 
            onClick={onClick}
        >
            {children}
        </div>
    )
}));

// Mock the MUI icons
vi.mock("@mui/icons-material/Movie", () => ({
    __esModule: true,
    default: () => <div data-testid="movie-icon">Movie Icon</div>
}));

vi.mock("@mui/icons-material/ExitToApp", () => ({
    __esModule: true,
    default: () => <div data-testid="exit-icon">Exit Icon</div>
}));

describe("Basic rendering", () => {
    const testOnLogout = vi.fn();

    beforeEach(() => {
        vi.clearAllMocks();
    });

    it("Should not render the logout section if not authenticated", () => {
        render(<Header isAuthenticated={false} onLogout={testOnLogout} />);

        // Check that the app bar and toolbar are rendered
        expect(screen.getByTestId("app-bar")).toBeInTheDocument();
        expect(screen.getByTestId("toolbar")).toBeInTheDocument();
        
        // Check that the logo is rendered
        expect(screen.getByTestId("movie-icon")).toBeInTheDocument();
        expect(screen.getByText("SkyFox Cinema")).toBeInTheDocument();
        
        // Check that the logout section is not rendered
        expect(screen.queryByText("Logout")).not.toBeInTheDocument();
        expect(screen.queryByTestId("exit-icon")).not.toBeInTheDocument();
    });

    it("Should render the logout section if authenticated", () => {
        render(<Header isAuthenticated={true} onLogout={testOnLogout} />);

        // Check that the app bar and toolbar are rendered
        expect(screen.getByTestId("app-bar")).toBeInTheDocument();
        expect(screen.getByTestId("toolbar")).toBeInTheDocument();
        
        // Check that the logo is rendered
        expect(screen.getByTestId("movie-icon")).toBeInTheDocument();
        expect(screen.getByText("SkyFox Cinema")).toBeInTheDocument();
        
        // Check that the logout section is rendered
        expect(screen.getByText("Logout")).toBeInTheDocument();
        expect(screen.getByTestId("exit-icon")).toBeInTheDocument();
        
        // Check that clicking on the logout box calls the onLogout function
        const logoutBox = screen.getAllByTestId("box")[1]; // Second Box is the logout section
        fireEvent.click(logoutBox);
        expect(testOnLogout).toHaveBeenCalledTimes(1);
    });
});
