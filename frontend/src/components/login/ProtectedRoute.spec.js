import React from "react";
import { render } from "@testing-library/react";
import ProtectedRoute from "./ProtectedRoute";
import * as reactRouterDom from "react-router-dom";

// Mock the react-router-dom hooks and components
vi.mock("react-router-dom", () => ({
    Navigate: vi.fn(({ to, state, replace }) => (
        <div data-testid="navigate" data-to={to} data-state={JSON.stringify(state)} data-replace={replace ? "true" : "false"}>
            Navigate Component
        </div>
    )),
    useLocation: vi.fn()
}));

describe("Basic Rendering", () => {
    const TestComponent = () => <div data-testid="test-component">Test Component</div>;
    const mockLocation = { pathname: "/current-path" };

    beforeEach(() => {
        // Reset mocks
        vi.clearAllMocks();
        
        // Setup location mock
        reactRouterDom.useLocation.mockReturnValue(mockLocation);
    });

    it("Should render children if authenticated", () => {
        const { getByTestId } = render(
            <ProtectedRoute isAuthenticated={true}>
                <TestComponent />
            </ProtectedRoute>
        );

        // Check that the children component is rendered
        expect(getByTestId("test-component")).toBeInTheDocument();
    });

    it("Should render Navigate if not authenticated", () => {
        const { getByTestId } = render(
            <ProtectedRoute isAuthenticated={false}>
                <TestComponent />
            </ProtectedRoute>
        );

        // Check that Navigate is rendered with correct props
        const navigateElement = getByTestId("navigate");
        expect(navigateElement).toBeInTheDocument();
        expect(navigateElement.getAttribute("data-to")).toBe("/login");
        
        // Check that the state contains the current location
        const stateAttr = navigateElement.getAttribute("data-state");
        const state = JSON.parse(stateAttr);
        expect(state.from).toEqual(mockLocation);
        
        // Check that replace is true
        expect(navigateElement.getAttribute("data-replace")).toBe("true");
    });
});
