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
import React, { FC, useState } from 'react';
import { HashRouter as Router, Link } from 'react-router-dom';

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

interface SidebarProps {
  open: boolean;
  handleDrawerClose(): void;
}

const Sidebar: FC<SidebarProps> = ({ handleDrawerClose, open }) => {
  const classes = useStyles();

  // Sidebar element
  const [testListOpen, setTestListOpen] = useState(false);
  const showTestList = () => {
    setTestListOpen(!testListOpen);
  };
  let colorTheme;
  let darkMode;

  React.useEffect(()=>{
  darkMode = localStorage.getItem('dark-mode');
  if(darkMode==='false') {
    colorTheme="white"
  }
  if(darkMode==='true') {
    colorTheme="black";
  }
  setCheckTestList(colorTheme);
  setCheckTestList1(colorTheme)
  setCheckTestList2(colorTheme)
  }, [darkMode]);
  const [checkTestList, setCheckTestList] = useState(colorTheme);
  const [checkTestList1, setCheckTestList1] = useState(colorTheme);
  const [checkTestList2, setCheckTestList2] = useState(colorTheme);
  const handleClick =  () => {
    setCheckTestList1(colorTheme)
    setCheckTestList2(colorTheme)
    setCheckTestList("#DCDCDC");
  }
  const handleClick1 =  () => {
    setCheckTestList(colorTheme);
    setCheckTestList2(colorTheme);
    setCheckTestList1("#DCDCDC");
  }
  const handleClick2 =  () => {
    setCheckTestList1(colorTheme);
    setCheckTestList(colorTheme);
    setCheckTestList2("#DCDCDC");
  }
  const menuItems = (
    <div>
      <ListItem button={true} component={Link} to="/">
        <ListItemIcon>
          <DashboardIcon />
        </ListItemIcon>
        <ListItemText primary="Dashboard" />
      </ListItem>
      <ListItem button={true} component={Link} to="/monitoring">
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
        {open ? (
          <ExpandLessIcon onClick={showTestList} />
        ) : (
          <ExpandMoreIcon onClick={showTestList} />
        )}
      </ListItem>
      {/* Nested List */}
      <Collapse in={testListOpen} timeout="auto" unmountOnExit={true}>
        <List component="div" disablePadding={true}>
          <ListItem
            button={true}
            component={Link}
            to="/ping"
            className={classes.nested}
            style={{backgroundColor: checkTestList}}
            onClick={handleClick}
          >
            <ListItemIcon>
              <DashboardIcon />
            </ListItemIcon>
            <ListItemText primary="Ping" />
          </ListItem>
          <ListItem
            button={true}
            component={Link}
            to="/floodping"
            className={classes.nested}
            style={{backgroundColor: checkTestList}}
            onClick={handleClick1}
          >
            <ListItemIcon>
              <DashboardIcon />
            </ListItemIcon>
            <ListItemText primary="FloodPing" />
          </ListItem>
          <ListItem
            button={true}
            component={Link}
            to="/jitter"
            className={classes.nested}
            style={{backgroundColor: checkTestList}}
            onClick={handleClick2}
          >
            <ListItemIcon>
              <DashboardIcon />
            </ListItemIcon>
            <ListItemText primary="Jitter" />
          </ListItem>
        </List>
      </Collapse>
      <ListItem button={true} component={Link} to="/configurations">
        <ListItemIcon>
          <SettingsIcon />
        </ListItemIcon>
        <ListItemText primary="Config" />
      </ListItem>
    </div>
  );

  return (
    <div className={classes.root}>
      <Router>
        <Drawer
          variant="permanent"
          classes={{
            paper: clsx(classes.drawerPaper, !open && classes.drawerPaperClose)
          }}
          open={open}
        >
          <div className={classes.toolbarIcon}>
            <IconButton onClick={handleDrawerClose}>
              <ChevronLeftIcon />
            </IconButton>
          </div>
          <Divider />
          <List>{menuItems}</List>
        </Drawer>
      </Router>
    </div>
  );
};

export default Sidebar;
