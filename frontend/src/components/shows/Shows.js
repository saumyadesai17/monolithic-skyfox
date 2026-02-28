import React, { useState } from "react";
import { useLocation, useNavigate } from "react-router-dom";
import {
    Avatar,
    Backdrop,
    Button,
    CircularProgress,
    List,
    ListItem,
    ListItemAvatar,
    ListItemText,
    Typography
} from "@mui/material";
import styles from "./styles/showsStyles"
import LocalMoviesIcon from "@mui/icons-material/LocalMovies";
import ArrowBackIcon from '@mui/icons-material/ArrowBack';
import ArrowForwardIcon from '@mui/icons-material/ArrowForward';
import useShows from "./hooks/useShows";
import { HEADER_DATE_FORMAT, INR_SYMBOL } from "../../Constants"
import { dateFromSearchString, nextDateLocation, previousDateLocation } from "./services/dateService";
import ShowsRevenue from "./ShowsRevenue";
import useShowsRevenue from "./hooks/useShowsRevenue";
import SeatSelectionDialog from "./SeatSelectionDialog";

const Shows = () => {
    const location = useLocation();
    const navigate = useNavigate();
    const classes = styles();

    const showsDate = dateFromSearchString(location.search);

    const {shows, showsLoading} = useShows(showsDate);
    const {showsRevenue, updateShowsRevenue, showsRevenueLoading} = useShowsRevenue(showsDate);
    const [showSelectSeatDialog, setShowSelectSeatDialog] = useState(false);
    const emptyShow = {
        "id": "",
        "date": "",
        "cost": "",
        "movie": {
            "id": "",
            "name": "",
            "duration": "",
            "plot": ""
        },
        "slot": {
            "id": "",
            "name": "",
            "startTime": "",
            "endTime": ""
        }
    };
    const [selectedShow, setSelectedShow] = useState(emptyShow);

    return (
        <>
            <div className={classes.cardHeader}>
                <Typography variant="h4" className={classes.showsHeader}>
                    Shows ({showsDate.format(HEADER_DATE_FORMAT)})
                </Typography>
                <ShowsRevenue showsRevenue={showsRevenue} showsRevenueLoading={showsRevenueLoading}/>
            </div>
            <List className={classes.listRoot}>
                {
                    shows.map(show => (
                        <div key={show.id} className={classes.showContainer}>
                            <ListItem style={{cursor: 'pointer'}} onClick={() => {
                                setSelectedShow(show);
                                setShowSelectSeatDialog(true);
                            }}>
                                <ListItemAvatar classes={{root: classes.localMoviesIcon}}>
                                    <Avatar>
                                        <LocalMoviesIcon/>
                                    </Avatar>
                                </ListItemAvatar>
                                <ListItemText primary={show.movie.name} secondary={
                                    <>
                                        <Typography
                                            component="span"
                                            variant="body2"
                                            className={classes.slotTime}
                                            color="textPrimary"
                                        >
                                            {show.slot.startTime}
                                        </Typography>
                                    </>
                                }/>
                                <ListItemText primary={`${INR_SYMBOL}${show.cost}`} className={classes.price}
                                              primaryTypographyProps={{variant: 'h6', color: 'secondary'}}
                                />
                            </ListItem>
                        </div>
                    ))
                }
            </List>

            <SeatSelectionDialog selectedShow={selectedShow} updateShowsRevenue={updateShowsRevenue}
                                 open={showSelectSeatDialog}
                                 onClose={() => setShowSelectSeatDialog(false)}/>

            <div className={classes.buttons}>
                <Button
                    onClick={() => {
                        navigate(previousDateLocation(location, showsDate));
                    }}
                    startIcon={<ArrowBackIcon/>}
                    color="primary"
                    className={classes.navigationButton}
                >
                    Previous Day
                </Button>
                <Button
                    onClick={() => {
                        navigate(nextDateLocation(location, showsDate));
                    }}
                    endIcon={<ArrowForwardIcon/>}
                    color="primary"
                    className={classes.navigationButton}
                >
                    Next Day
                </Button>
            </div>
            <Backdrop className={classes.backdrop} open={showsLoading}>
                <CircularProgress color="inherit"/>
            </Backdrop>
        </>
    );
};

export default Shows;
