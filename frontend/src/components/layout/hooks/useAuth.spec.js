import {act, renderHook} from "@testing-library/react-hooks";
import useAuth from "./useAuth";
import {isLoggedIn, login, logout} from "../../../helpers/authService";

vi.mock("../../../helpers/authService", () => ({
    __esModule: true,
    isLoggedIn: vi.fn(),
    login: vi.fn(),
    logout: vi.fn()
}));

describe("Basic logic", () => {

    it("should respect user's logged in status initially", () => {
        isLoggedIn.mockReturnValue(true);
        const renderHookResult = renderHook(() => useAuth());
        const {result, waitForNextUpdate} = renderHookResult;

        waitForNextUpdate();

        const {isAuthenticated} = result.current;
        expect(isAuthenticated).toBe(true);
    });

    it("should login successfuly", async () => {
        const testUsername = "testUsername";
        const testPassword = "testPassword";
        isLoggedIn.mockReturnValue(true);
        
        // Mock the login function to return userDetails when called with specific credentials
        login.mockImplementation((username, password) => {
            if (username === testUsername && password === testPassword) {
                return Promise.resolve("userDetails");
            }
            return Promise.reject(new Error("Unexpected credentials"));
        });
        
        const renderHookResult = renderHook(() => useAuth());
        const {result} = renderHookResult;

        const {handleLogin} = result.current;

        let userDetails;
        await act(async () => {
            userDetails = await handleLogin(testUsername, testPassword);
        });

        const {isAuthenticated} = result.current;
        // noinspection JSUnusedAssignment
        expect(userDetails).toBe("userDetails");
        expect(isAuthenticated).toBe(true);
    });

    it("should not login if not authenticated", async () => {
        const testUsername = "testUsername";
        const testPassword = "testPassword";
        isLoggedIn.mockReturnValue(false);
        
        // Mock the login function to reject when called with specific credentials
        login.mockImplementation((username, password) => {
            if (username === testUsername && password === testPassword) {
                return Promise.reject("unused");
            }
            return Promise.reject(new Error("Unexpected credentials"));
        });
        
        const renderHookResult = renderHook(() => useAuth());
        const {result} = renderHookResult;

        const {handleLogin} = result.current;

        try {
            await act(async () => {
                await handleLogin(testUsername, testPassword);
                fail("Should not authenticate");
            });
        } catch (e) {
            const {isAuthenticated} = result.current;
            expect(isAuthenticated).toBe(false);
        }
    });

    it("should logout successfuly", () => {
        isLoggedIn.mockReturnValue(true);
        const renderHookResult = renderHook(() => useAuth());
        const {result} = renderHookResult;

        const {handleLogout} = result.current;
        act(() => {
            handleLogout();
        });

        const {isAuthenticated} = result.current;
        expect(logout).toBeCalledTimes(1);
        expect(isAuthenticated).toBe(false);
    });
});
