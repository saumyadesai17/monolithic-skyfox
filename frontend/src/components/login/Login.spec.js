import React from "react";
import { render, screen } from "@testing-library/react";
import Login from "./Login";
import useLogin from "./hooks/useLogin";
import * as reactRouterDom from "react-router-dom";

// Mock the react-router-dom hooks
vi.mock("react-router-dom", () => ({
    useLocation: vi.fn(),
    useNavigate: vi.fn()
}));

vi.mock("./hooks/useLogin", () => ({
    __esModule: true,
    default: vi.fn()
}));

vi.mock("./services/loginFormService", () => ({
    __esModule: true,
    initialValues: "initialValues",
    formSchema: "formSchema"
}));

// Mock the Formik component and useField
vi.mock("formik", () => ({
    Formik: ({ children, initialValues, validationSchema, onSubmit }) => (
        <div data-testid="formik" data-initialvalues={initialValues} data-validationschema={validationSchema}>
            {typeof children === 'function' ? children({ isValid: true }) : children}
        </div>
    ),
    Form: ({ children }) => <form data-testid="form">{children}</form>,
    useField: vi.fn().mockReturnValue([
        { value: '', onChange: vi.fn(), onBlur: vi.fn() },
        { touched: false, error: '' }
    ])
}));

describe("Basic Rendering", () => {
    const testOnLogin = vi.fn();
    const testHandleLogin = vi.fn();
    const testFrom = "/testFrom";
    const TestErrorComponent = () => <div data-testid="error-component">Error Component</div>;
    const mockNavigate = vi.fn();

    beforeEach(() => {
        // Reset mocks
        vi.clearAllMocks();
        
        // Setup mocks
        reactRouterDom.useNavigate.mockReturnValue(mockNavigate);
        reactRouterDom.useLocation.mockReturnValue({
            state: { from: { pathname: testFrom } }
        });
        
        useLogin.mockReturnValue({
            errorMessage: () => <TestErrorComponent />,
            handleLogin: testHandleLogin
        });
    });

    it("should navigate to from url when authenticated", () => {
        render(<Login isAuthenticated={true} onLogin={testOnLogin} />);

        // Check that navigate was called with the correct arguments
        expect(mockNavigate).toHaveBeenCalledWith(testFrom, { replace: true });
    });

    it("should render login form when not authenticated", () => {
        render(<Login isAuthenticated={false} onLogin={testOnLogin} />);

        // Check that the form is rendered
        expect(screen.getByTestId("form")).toBeInTheDocument();
        
        // Check that the error component is rendered
        expect(screen.getByTestId("error-component")).toBeInTheDocument();
        
        // Check Formik props
        const formikElement = screen.getByTestId("formik");
        expect(formikElement.getAttribute("data-initialvalues")).toBe("initialValues");
        expect(formikElement.getAttribute("data-validationschema")).toBe("formSchema");
    });
});
