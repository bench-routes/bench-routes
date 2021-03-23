import { makeStyles } from '@material-ui/core';

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
  },
  popupTitle: {
    padding: '2rem 2rem 0 2rem',
    borderTop: '0.3rem solid #2195f1'
  },
  popupContent: {
    padding: '1rem 2rem 2rem 2rem'
  },
  popupButton: {
    '&:focus': {
      outline: 'none'
    },
    fontWeight: 800,
    fontSize: '0.8rem'
  },
  popupActions: {
    backgroundColor: '#f2f2f5',
    padding: '1rem 2rem'
  }
}));

export default useStyles;
