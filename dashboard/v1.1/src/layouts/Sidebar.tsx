import {
  Collapse,
  Divider,
  Drawer,
  IconButton,
  List,
  ListItem,
  ListItemIcon,
  ListItemText
} from '@material-ui/core';
import { makeStyles } from '@material-ui/core/styles';
import {
  AccessAlarm as AccessAlarmIcon,
  ChevronLeft as ChevronLeftIcon,
  Dashboard as DashboardIcon,
  ExpandLess as ExpandLessIcon,
  ExpandMore as ExpandMoreIcon,
  NetworkCheck as NetworkCheckIcon,
  Settings as SettingsIcon
} from '@material-ui/icons';
import clsx from 'clsx';
import React, { useState } from 'react';

const drawerWidth = 240;

const useStyles = makeStyles(theme => ({
  // root class
  root: {
    display: 'flex'
  },
  // Nested lists
  nested: {
    paddingLeft: theme.spacing(4),
    backgroundColor: '#DCDCDC'
  },
  // Drawer styles
  drawerPaper: {
    position: 'relative',
    whiteSpace: 'nowrap',
    width: drawerWidth,
    transition: theme.transitions.create('width', {
      easing: theme.transitions.easing.sharp,
      duration: theme.transitions.duration.enteringScreen
    })
  },
  drawerPaperClose: {
    overflowX: 'hidden',
    transition: theme.transitions.create('width', {
      easing: theme.transitions.easing.sharp,
      duration: theme.transitions.duration.leavingScreen
    }),
    width: theme.spacing(7),
    [theme.breakpoints.up('sm')]: {
      width: theme.spacing(9)
    }
  },
  toolbarIcon: {
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'flex-end',
    padding: '0 8px',
    ...theme.mixins.toolbar
  },
  // SettingsIcon
  settings: {
    position: 'absolute',
    bottom: 0,
    alignItems: 'center'
  }
}));

export default function Sidebar(props) {
  const classes = useStyles();

  // Sidebar element
  const [testListOpen, setTestListOpen] = useState(false);
  const showTestList = () => {
    setTestListOpen(!testListOpen);
  };
  const menuItems = (
    <div>
      <ListItem button={true}>
        <ListItemIcon>
          <DashboardIcon />
        </ListItemIcon>
        <ListItemText primary="Dashboard" />
      </ListItem>
      <ListItem button={true}>
        <ListItemIcon>
          <AccessAlarmIcon />
        </ListItemIcon>
        <ListItemText primary="Monitoring" />
      </ListItem>
      <ListItem button={true}>
        <ListItemIcon>
          <NetworkCheckIcon />
        </ListItemIcon>
        <ListItemText primary="Tests" onClick={showTestList} />
        {props.open ? (
          <ExpandLessIcon onClick={showTestList} />
        ) : (
          <ExpandMoreIcon onClick={showTestList} />
        )}
      </ListItem>
      {/* Nested List */}
      <Collapse in={testListOpen} timeout="auto" unmountOnExit={true}>
        <List component="div" disablePadding={true}>
          <ListItem button={true} className={classes.nested}>
            <ListItemIcon>
              <DashboardIcon />
            </ListItemIcon>
            <ListItemText primary="Ping" />
          </ListItem>
          <ListItem button={true} className={classes.nested}>
            <ListItemIcon>
              <DashboardIcon />
            </ListItemIcon>
            <ListItemText primary="FloodPing" />
          </ListItem>
          <ListItem button={true} className={classes.nested}>
            <ListItemIcon>
              <DashboardIcon />
            </ListItemIcon>
            <ListItemText primary="Jitter" />
          </ListItem>
        </List>
      </Collapse>
    </div>
  );

  return (
    <div className={classes.root}>
      <Drawer
        variant="permanent"
        classes={{
          paper: clsx(
            classes.drawerPaper,
            !props.open && classes.drawerPaperClose
          )
        }}
        open={props.open}
      >
        <div className={classes.toolbarIcon}>
          <IconButton onClick={props.handleDrawerClose}>
            <ChevronLeftIcon />
          </IconButton>
        </div>
        <Divider />
        <List>{menuItems}</List>
        <ListItem button={true} className={classes.settings}>
          <ListItemIcon>
            <SettingsIcon />
          </ListItemIcon>
          <ListItemText primary="Settings" />
        </ListItem>
      </Drawer>
    </div>
  );
}
