import {renderHook} from "@testing-library/react-hooks";
import useShowsRevenue from "./useShowsRevenue";
import showsService from "../services/showsService";
import moment from "moment";

vi.mock("../services/showsService", () => ({
    __esModule: true,
    default: {
        getRevenue: vi.fn()
    }
}));

describe("Basic logic", () => {
    let showDate;

    beforeEach(() => {
        showDate = moment("2020-01-01", "YYYY-MM-DD");
        
        // Mock the getRevenue function to return revenue when called with a specific date
        showsService.getRevenue.mockImplementation((date) => {
            if (date === "2020-01-01") {
                return Promise.resolve(549.99);
            }
            return Promise.reject(new Error("Unexpected date"));
        });
    });

    it("Should initialize the hook with zero shows revenue and loading", () => {
        const {result} = renderHook(() => useShowsRevenue(showDate));

        const {showsRevenue, showsRevenueLoading} = result.current;

        expect(showsRevenue).toEqual(0);
        expect(showsRevenueLoading).toBe(true);
    });

    it("Should get shows revenue and finish loading after mount", async () => {
        const {result, waitForNextUpdate} = renderHook(() => useShowsRevenue(showDate));

        await waitForNextUpdate();
        const {showsRevenue, showsRevenueLoading} = result.current;

        expect(showsRevenue).toEqual(549.99);
        expect(showsRevenueLoading).toBe(false);
    });
});
