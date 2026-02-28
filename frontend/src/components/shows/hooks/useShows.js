import {useEffect, useState} from 'react';
import showsService from "../services/showsService";
import {QUERY_DATE_FORMAT} from "../../../Constants";

const useShows = (showsDate) => {
    const [showsLoading, setShowsLoading] = useState(true);
    const [shows, setShows] = useState([]);

    useEffect(() => {
        const formattedDate = showsDate.format(QUERY_DATE_FORMAT);
        showsService.fetchAll(formattedDate).then(shows => {
            setShowsLoading(false);
            setShows(shows);
        });
        // eslint-disable-next-line
    }, []);

    return {
        shows: shows,
        showsLoading: showsLoading
    };
}

export default useShows;
