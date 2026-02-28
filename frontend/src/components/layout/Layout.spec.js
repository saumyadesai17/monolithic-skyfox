import React from "react";
import { render } from "@testing-library/react";
import Layout from "./Layout";
import useAuth from "./hooks/useAuth";
import Header from "../header/Header";
import RootRouter from "../router/RootRouter";

const testHandleLogin = vi.fn();
const testHandleLogout = vi.fn();

// Mock the child components
vi.mock("../header/Header", () => ({
    __esModule: true,
    default: vi.fn(() => <div data-testid="header">Header Component</div>)
}));

vi.mock("../router/RootRouter", () => ({
    __esModule: true,
    default: vi.fn(() => <div data-testid="root-router">Root Router Component</div>)
}));

vi.mock("./hooks/useAuth", () => ({
    __esModule: true,
    default: vi.fn()
}));

describe('Basic rendering', function () {
    beforeEach(() => {
        // Clear mock calls before each test
        Header.mockClear();
        RootRouter.mockClear();
    });

    it("Should render correctly", () => {
        // Setup auth hook mock
        useAuth.mockReturnValue({
            isAuthenticated: true,
            handleLogin: testHandleLogin,
            handleLogout: testHandleLogout
        });
        
        // Render the component
        render(<Layout/>);
        
        // Check Header props
        expect(Header).toHaveBeenCalledWith(
            expect.objectContaining({
                onLogout: testHandleLogout,
                isAuthenticated: true
            }),
            expect.anything()
        );
        
        // Check RootRouter props
        expect(RootRouter).toHaveBeenCalledWith(
            expect.objectContaining({
                onLogin: testHandleLogin,
                isAuthenticated: true
            }),
            expect.anything()
        );
    });
});
