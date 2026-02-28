import { makeStyles } from "@mui/styles";

export default makeStyles((theme) =>
    ({
        errorContent: {
            display: 'flex',
            justifyContent: 'center',
            alignItems: 'center',
            flexDirection: 'column',
            paddingBottom: '25px'
        },
        errorIcon: {
            color: theme.palette.error.main,
            height: '400px',
            width: '400px',
            opacity: '0.75'
        }
    })
);
