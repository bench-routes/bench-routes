import { Fab, makeStyles, Tooltip } from '@material-ui/core';
import { PostAdd as PostAddIcon } from '@material-ui/icons';
import React from 'react';
import { Link } from 'react-router-dom';

const useStyles = makeStyles(theme => ({
  inputFab: {
    position: 'fixed',
    right: '2rem',
    bottom: '2rem'
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
        </Fab>
      </Tooltip>
    </Link>
  );
};

export default QuickInputFab;
