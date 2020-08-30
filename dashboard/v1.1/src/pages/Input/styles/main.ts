import { makeStyles } from "@material-ui/core";

const useStyles = makeStyles(theme => ({
  pageheader: {
    marginBottom: '20px'
  },
  container: {
    border: '0.1rem solid #e3e3e3',
    marginTop: '1rem'
  },
  controls: {
    marginBottom: '1.4rem',
    marginTop: '1.4rem'
  },
  additionalParams: {
    marginBottom: '1.4rem'
  },
  btn: {
    width: '100%',
    minHeight: '100%',
    fontSize: '16px'
  },
  marginRtSm: {
    marginRight: '0.2rem'
  },
  marginTopMd: {
    marginTop: '1rem'
  },
  params: {
    margin: '2%'
  }
}));

export default useStyles;
