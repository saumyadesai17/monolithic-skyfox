import {renderHook} from "@testing-library/react-hooks";
import useShows from "./useShows";
import showsService from "../services/showsService";
import moment from "moment";

vi.mock("../services/showsService", () => ({
    __esModule: true,
    default: {
        fetchAll: vi.fn()
    }
}));

describe("Basic logic", () => {
    let showDate;

    beforeEach(() => {
        showDate = moment("2020-01-01", "YYYY-MM-DD");
        
        // Mock the fetchAll function to return shows when called with a specific date
        showsService.fetchAll.mockImplementation((date) => {
            if (date === "2020-01-01") {
                return Promise.resolve(["show 1", "show 2"]);
            }
            return Promise.reject(new Error("Unexpected date"));
        });
    });

    it("Should initialize the hook with empty shows and loading", () => {
        const {result} = renderHook(() => useShows(showDate));

        const {shows, showsLoading} = result.current;

        expect(shows).toEqual([]);
        expect(showsLoading).toBe(true);
    });

    it("Should get shows and finish loading after mount", async () => {
        const {result, waitForNextUpdate} = renderHook(() => useShows(showDate));

        await waitForNextUpdate();
        const {shows, showsLoading} = result.current;

        expect(shows).toEqual(["show 1", "show 2"]);
        expect(showsLoading).toBe(false);
    });
});
