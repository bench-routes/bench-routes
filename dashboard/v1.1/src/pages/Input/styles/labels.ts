import { fade, makeStyles } from "@material-ui/core";

const useStyles = makeStyles(theme => ({
  root: {
    fontSize: 13,
    width: '100%'
  },
  button: {
    fontSize: 13,
    width: '100%',
    textAlign: 'left',
    paddingBottom: 5,
    color: '#586069',
    fontWeight: 600,
    '&:hover,&:focus': {
      color: '#0366d6'
    },
    '& span': {
      width: '100%'
    },
    '& svg': {
      width: 16,
      height: 16
    }
  },
  tag: {
    height: 20,
    padding: '.15em 4px',
    fontWeight: 600,
    lineHeight: '15px',
    borderRadius: 2,
    minWidth: 'max-content',
    marginRight: '0.2rem'
  },
  popper: {
    border: '1px solid rgba(27,31,35,.15)',
    boxShadow: '0 3px 12px rgba(27,31,35,.15)',
    borderRadius: 3,
    width: 300,
    zIndex: theme.zIndex.modal,
    fontSize: 13,
    color: '#586069',
    backgroundColor: '#f6f8fa'
  },
  header: {
    borderBottom: '1px solid #e1e4e8',
    padding: '8px 10px',
    fontWeight: 600,
    display: 'flex',
    justifyContent: 'space-between'
  },
  inputBase: {
    padding: 10,
    width: '100%',
    borderBottom: '1px solid #dfe2e5',
    '& input': {
      borderRadius: 4,
      backgroundColor: theme.palette.common.white,
      padding: 8,
      transition: theme.transitions.create(['border-color', 'box-shadow']),
      border: '1px solid #ced4da',
      fontSize: 14,
      '&:focus': {
        boxShadow: `${fade(theme.palette.primary.main, 0.25)} 0 0 0 0.2rem`,
        borderColor: theme.palette.primary.main
      }
    }
  },
  paper: {
    boxShadow: 'none',
    margin: 0,
    color: '#586069',
    fontSize: 13
  },
  option: {
    minHeight: 'auto',
    alignItems: 'flex-start',
    padding: 8,
    '&[aria-selected="true"]': {
      backgroundColor: 'transparent'
    },
    '&[data-focus="true"]': {
      backgroundColor: theme.palette.action.hover
    }
  },
  popperDisablePortal: {
    position: 'relative'
  },
  iconSelected: {
    width: 17,
    height: 17,
    marginRight: 5,
    marginLeft: -2
  },
  color: {
    width: 14,
    height: 14,
    flexShrink: 0,
    borderRadius: 20,
    marginRight: '1rem',
    marginTop: 2
  },
  text: {
    flexGrow: 1
  },
  close: {
    opacity: 0.6,
    width: 18,
    height: 18
  },
  noOption: {
    padding: '0.2rem',
    display: 'flex',
    cursor: 'pointer',
    '&:hover': {
      textDecoration: 'underline',
      textDecorationColor: '#1976D2'
    }
  },
  value: {
    overflowX: 'scroll',
    display: 'flex'
  }
}));

export default useStyles;
