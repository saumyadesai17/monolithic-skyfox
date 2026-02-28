import {useEffect, useState} from "react";
import {QUERY_DATE_FORMAT} from "../../../Constants";
import showsService from "../services/showsService";

const useShowsRevenue = (showsDate) => {
    const [showsRevenueLoading, setShowsRevenueLoading] = useState(true);
    const [showsRevenue, setShowsRevenue] = useState(0);

    useEffect(() => {
        updateShowsRevenue();
        // eslint-disable-next-line
    }, []);

    const updateShowsRevenue = () => {
        const formattedDate = showsDate.format(QUERY_DATE_FORMAT);

        showsService.getRevenue(formattedDate).then(showsRevenue => {
            setShowsRevenueLoading(false);
            setShowsRevenue(showsRevenue);
        });
    };

    return {
        showsRevenue: showsRevenue,
        updateShowsRevenue: updateShowsRevenue,
        showsRevenueLoading: showsRevenueLoading
    };
};

export default useShowsRevenue;
