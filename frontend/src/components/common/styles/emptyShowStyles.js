import { makeStyles } from "@mui/styles";

export default makeStyles((theme) =>
    ({
        emptyShowsLayout: {
            display: "flex",
            flexDirection: "column",
            justifyContent: "flex-start",
            alignItems: "flex-start",
            padding: "30px"
        },
        emptyShowsIcon: {
            width: "75px",
            height: "75px"
        },
        emptyShowsContainer: {
            opacity: "35%",
            display: "flex",
            flexDirection: "column",
            justifyContent: "center",
            alignItems: "center",
        }
    })
);
