import {act, renderHook} from "@testing-library/react-hooks";
import useLogin from "./useLogin";
import React from "react";
import { render, screen } from "@testing-library/react";

describe("Basic logic", () => {

    const testUsername = "testUsername";
    const testPassword = "testPassword";
    const loginValues = {
        username: testUsername,
        password: testPassword
    };

    it("should initially not show error message", () => {
        const testOnLogin = vi.fn();
        const renderHookResult = renderHook(() => useLogin(testOnLogin));
        const result = renderHookResult.result;
        const {errorMessage} = result.current;

        expect(errorMessage()).toBe(undefined);
    });

    it("should not show error message if logged in succesfully", async () => {
        const testOnLogin = vi.fn().mockImplementation((username, password) => {
            if (username === testUsername && password === testPassword) {
                return Promise.resolve("Unused");
            }
            return Promise.reject(new Error("Unexpected arguments"));
        });
        
        const renderHookResult = renderHook(() => useLogin(testOnLogin));
        const result = renderHookResult.result;
        const {handleLogin} = result.current;

        await act(() => handleLogin(loginValues));

        const {errorMessage} = result.current;
        expect(testOnLogin).toBeCalledTimes(1);
        expect(testOnLogin).toHaveBeenCalledWith(testUsername, testPassword);
        expect(errorMessage()).toBe(undefined);
    });

    it("should show error message if 401 returned", async () => {
        const testOnLogin = vi.fn().mockImplementation((username, password) => {
            if (username === testUsername && password === testPassword) {
                return Promise.reject({
                    response: {
                        status: 401
                    }
                });
            }
            return Promise.reject(new Error("Unexpected arguments"));
        });
        
        const renderHookResult = renderHook(() => useLogin(testOnLogin));
        const result = renderHookResult.result;
        const {handleLogin} = result.current;

        await act(() => handleLogin(loginValues));

        const {errorMessage} = result.current;
        
        // Render the error message component and check its content
        const { container } = render(errorMessage());
        expect(testOnLogin).toBeCalledTimes(1);
        expect(testOnLogin).toHaveBeenCalledWith(testUsername, testPassword);
        expect(container.textContent).toBe("Login failed");
    });

    it("should not show error message if non-401 error", async () => {
        const testError = "test error";
        const testOnLogin = vi.fn().mockImplementation((username, password) => {
            if (username === testUsername && password === testPassword) {
                return Promise.reject(testError);
            }
            return Promise.reject(new Error("Unexpected arguments"));
        });
        
        const renderHookResult = renderHook(() => useLogin(testOnLogin));
        const result = renderHookResult.result;
        const {handleLogin} = result.current;

        try {
            await act(() => handleLogin(loginValues));
            fail("Error not rethrown");
        } catch (err) {
            const {errorMessage} = result.current;
            expect(testOnLogin).toBeCalledTimes(1);
            expect(testOnLogin).toHaveBeenCalledWith(testUsername, testPassword);
            expect(errorMessage()).toBe(undefined);
            expect(err).toBe(testError)
        }
    });
});
