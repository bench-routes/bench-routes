import {
  AppBar,
  Badge,
  IconButton,
  Toolbar,
  Typography
} from '@material-ui/core';
import { makeStyles } from '@material-ui/core/styles';
import {
  Menu as MenuIcon,
  Notifications as NotificationsIcon
} from '@material-ui/icons';
import clsx from 'clsx';
import React, { useState, FC } from 'react';
import LinearProgress from '@material-ui/core/LinearProgress';

const drawerWidth = 240;

const useStyles = makeStyles(theme => ({
  // AppBar styles
  appBar: {
    zIndex: theme.zIndex.drawer + 1,
    transition: theme.transitions.create(['width', 'margin'], {
      easing: theme.transitions.easing.sharp,
      duration: theme.transitions.duration.leavingScreen
    })
  },
  appBarShift: {
    marginLeft: drawerWidth,
    width: `calc(100% - ${drawerWidth}px)`,
    transition: theme.transitions.create(['width', 'margin'], {
      easing: theme.transitions.easing.sharp,
      duration: theme.transitions.duration.enteringScreen
    })
  },
  title: {
    flexGrow: 1
  },
  // Toolbar styles
  toolbar: {
    paddingRight: 24 // keep right padding when drawer closed
  },
  // IconMenu styles
  menuButton: {
    marginRight: 36
  },
  menuButtonHidden: {
    display: 'none'
  }
}));

interface NavbarProps {
  handleDrawerOpen(): void;
  open: boolean;
  getLoaderStatus(): boolean;
}

const Navbar: FC<NavbarProps> = ({
  handleDrawerOpen,
  open,
  getLoaderStatus
}) => {
  const classes = useStyles();
  const [showLoader, setShowLoader] = useState<boolean>(getLoaderStatus());

  return (
    <div>
      <AppBar
        position="absolute"
        className={clsx(classes.appBar, open && classes.appBarShift)}
      >
        <Toolbar className={classes.toolbar}>
          <IconButton
            edge="start"
            color="inherit"
            aria-label="open drawer"
            onClick={handleDrawerOpen}
            className={clsx(
              classes.menuButton,
              open && classes.menuButtonHidden
            )}
          >
            <MenuIcon />
          </IconButton>
          <Typography
            component="h1"
            variant="h6"
            color="inherit"
            noWrap={true}
            className={classes.title}
          >
            Bench Routes
          </Typography>
          <IconButton color="inherit">
            <Badge badgeContent={4} color="secondary">
              <NotificationsIcon />
            </Badge>
          </IconButton>
        </Toolbar>
        {showLoader ? <LinearProgress /> : null}
      </AppBar>
    </div>
  );
};

export default Navbar;
