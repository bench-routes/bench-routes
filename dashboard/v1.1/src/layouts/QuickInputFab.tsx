import { Fab, makeStyles, Tooltip } from '@material-ui/core';
import { PostAdd as PostAddIcon } from '@material-ui/icons';
import React from 'react';
import { Link } from 'react-router-dom';

const useStyles = makeStyles(theme => ({
  inputFab: {
    position: 'fixed',
    right: 0,
    bottom: '10rem',
    borderRadius: '10px 0 0 10px',
    width: '10rem',
    height: '3rem',
    zIndex: 12,
    transform: 'translateX(70%)',
    transition: '0.3s all ease-in'
  },
  tag: {
    fontWeight: 600,
    padding: 15
  }
}));

// Floating Action Button for Quick Input
const QuickInputFab = () => {
  const classes = useStyles();
  return (
    <Link to="/quick-input">
      <Tooltip placement="top" title="Quick Route Input">
        <Fab
          className={classes.inputFab}
          color="primary"
          aria-label="quick-input"
        >
          <PostAddIcon />
          <span className={classes.tag}>Quick Input</span>
        </Fab>
      </Tooltip>
    </Link>
  );
};

export default QuickInputFab;
