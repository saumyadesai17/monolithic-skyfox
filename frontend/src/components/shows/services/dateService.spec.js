import { describe, it, expect, beforeEach } from 'vitest';
import {dateFromSearchString, nextDateLocation, previousDateLocation} from "./dateService";
import moment from "moment";

const expectTodayOrYesterday = (dateFormatted) => {
    const todayFormatted = moment().format("YYYY-MM-DD");
    const yesterdayFormatted = moment().subtract(1, "days").format("YYYY-MM-DD");

    expect(dateFormatted === todayFormatted || dateFormatted === yesterdayFormatted).toBe(true);
};

describe("Basic Logic", () => {
    let testLocation;

    beforeEach(() => {
        testLocation = {
            path: "testPath",
            search: "testSearch"
        }
    });

    // Using Vitest's test.each instead of jest-each
    it.each([undefined, null, "", "invalid", "?invalid=2020-02-02", "?date=invalid", "?date=2020-45-45"])
        ("Should return current date if search is \"%s\"", (searchString) => {
            const actualDateFormatted = dateFromSearchString(searchString).format("YYYY-MM-DD");

            expectTodayOrYesterday(actualDateFormatted);
        });

    it("Should get valid date passed from search string", () => {
        const actualDateFormatted = dateFromSearchString("?date=2020-01-01").format("YYYY-MM-DD");

        expect(actualDateFormatted).toBe("2020-01-01");
    });

    it("Should get next date formatted", () => {
        const date = moment("2020-01-01", "YYYY-MM-DD");

        const actualNextDateFormatted = nextDateLocation(testLocation, date);

        expect(actualNextDateFormatted).toEqual({
            path: "testPath",
            search: "?date=2020-01-02"
        });
    });

    it("Should get previous date formatted", () => {
        const date = moment("2020-01-01", "YYYY-MM-DD");

        const actualNextDateFormatted = previousDateLocation(testLocation, date);

        expect(actualNextDateFormatted).toEqual({
            path: "testPath",
            search: "?date=2019-12-31"
        });
    });
});
