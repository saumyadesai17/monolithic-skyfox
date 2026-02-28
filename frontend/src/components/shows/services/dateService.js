import queryString from "query-string";
import moment from "moment";
import {QUERY_DATE_FORMAT} from "../../../Constants";

export const dateFromSearchString = (searchString) => {
    if (searchString && ("date" in queryString.parse(searchString))) {
        const parsedQuery = queryString.parse(searchString);

        if (moment(parsedQuery.date, QUERY_DATE_FORMAT).isValid()) {
            return moment(queryString.parse(searchString, QUERY_DATE_FORMAT).date);
        }
    }

    return moment();
};

export const nextDateLocation = (currentLocation, date) => {
    const nextDateLocationFormatted = date.add(1, 'days').format(QUERY_DATE_FORMAT);

    return locationWithFormattedDate(currentLocation, nextDateLocationFormatted)
};

export const previousDateLocation = (currentLocation, date) => {
    const previousDateLocationFormatted = date.subtract(1, 'days').format(QUERY_DATE_FORMAT);

    return locationWithFormattedDate(currentLocation, previousDateLocationFormatted)
};

const locationWithFormattedDate = (currentLocation, previousDateLocationFormatted) => {
    return {
        ...currentLocation,
        search: `?date=${previousDateLocationFormatted}`
    };
}
