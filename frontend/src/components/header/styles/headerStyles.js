// Define styles as an object with sx props
const headerStyles = {
  headerLink: {
    color: 'primary.contrastText',
    display: 'flex',
    justifyContent: 'flex-start',
    textDecoration: 'none'
  },
  logoutLink: {
    display: 'flex',
    justifyContent: 'flex-start',
    alignItems: 'center',
    cursor: 'pointer'
  },
  cinemaLogoIcon: {
    fontSize: '2.25em'
  },
  headerLogo: {
    marginLeft: '0.15em'
  },
  toolbar: {
    display: 'flex',
    justifyContent: 'space-between',
    padding: '0 4em'
  },
  cartButton: {
    color: 'white',
    padding: '10px'
  }
};

export default headerStyles;
