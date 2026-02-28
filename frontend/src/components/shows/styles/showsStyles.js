import { makeStyles } from "@mui/styles";

export default makeStyles((theme) =>
    ({
        cardHeader: {
            display: "flex",
            justifyContent: "space-between"
        },
        showContainer: {
            "& :hover": {
                backgroundColor: "#f9f8fd",
            }
        },
        localMoviesIcon: {
            "& :hover": {
                backgroundColor: "#bdbdbd",
            }
        },
        showsHeader: {
            padding: "15px 0 0 15px",
            display: "flex",
            fontWeight: "bold",
            alignSelf: "center"
        },
        backdrop: {
            zIndex: theme.zIndex?.drawer ? theme.zIndex.drawer + 1 : 1200, // Use default value if drawer is undefined
            color: '#fff',
        },
        listRoot: {
            width: '100%',
            backgroundColor: theme.palette?.background?.paper || '#fff' // Use default white if background.paper is undefined
        },
        price: {
            display: 'flex',
            justifyContent: 'flex-end',
        },
        slotTime: {
            color: theme.palette?.primary?.main || '#556cd6', // Use default primary color if primary.main is undefined
            fontWeight: "bold"
        },
        buttons: {
            display: "flex",
            justifyContent: 'space-between'
        },
        navigationButton: {
            margin: "20px"
        },
        paper: {
            width: '200',
            height: '500',
            backgroundColor: theme.palette?.background?.paper || '#fff', // Use default white if background.paper is undefined
            boxShadow: theme.shadows?.[5] || '0px 3px 5px -1px rgba(0,0,0,0.2)', // Default shadow if theme.shadows is undefined
            padding: theme.spacing ? theme.spacing(2, 4, 3) : '16px 32px 24px', // Default padding if theme.spacing is undefined
        }
    })
);
