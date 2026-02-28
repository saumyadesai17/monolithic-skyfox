import { makeStyles } from "@mui/styles";

export default makeStyles((theme) =>
    ({
        showsRevenueContainer: {
            padding: "15px 15px 0px 0px",
            display: "flex",
            alignSelf: "center",
            color: "green",
            fontWeight: "bold"
        },
        showsRevenueLoadingSpinner: {
            padding: "15px"
        }
    })
);
