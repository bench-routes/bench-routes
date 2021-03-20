import React, { FC } from 'react';
import { chartData, RouteDetails } from '../../utils/queryTypes';
import { formatTime } from '../../utils/brt';
import ChartComponent from './Chart';
import {
  JITTER_OPTIONS,
  PING_OPTIONS,
  DELAY_OPTIONS,
  RESPONSE_LENGTH_OPTIONS
} from '../../utils/constants/chart';

import AppBar from '@material-ui/core/AppBar';
import Tabs from '@material-ui/core/Tabs';
import Tab from '@material-ui/core/Tab';
import PropTypes from 'prop-types';
import Box from '@material-ui/core/Box';
import ArrowBackIcon from '@material-ui/icons/ArrowBack';
import { truncate } from '../../utils/stringManipulations';

interface RouteDetailsProps {
  routesChains: RouteDetails;
  showDetails(status: boolean, details: RouteDetails): void;
}

function a11yProps(index) {
  return {
    id: `simple-tab-${index}`,
    'aria-controls': `simple-tabpanel-${index}`
  };
}

function TabPanel(props) {
  const { children, value, index, ...other } = props;

  return (
    <div
      role="tabpanel"
      hidden={value !== index}
      id={`simple-tabpanel-${index}`}
      aria-labelledby={`simple-tab-${index}`}
      {...other}
    >
      {value === index && <Box p={3}>{children}</Box>}
    </div>
  );
}

TabPanel.propTypes = {
  children: PropTypes.node,
  index: PropTypes.any.isRequired,
  value: PropTypes.any.isRequired
};

const format = (data: RouteDetails) => {
  const responseDetailsDelay: chartData[] = [];
  const responseDetailsResponse: chartData[] = [];
  const pingMin: chartData[] = [];
  const pingMean: chartData[] = [];
  const pingMax: chartData[] = [];
  const jitter: chartData[] = [];
  const name = data.name;

  if (data.monitor.values) {
    for (const value of data.monitor.values) {
      responseDetailsDelay.push({
        y: value.value.delay,
        x: formatTime(value.timestamp)
      });
      responseDetailsResponse.push({
        y: value.value.resLength,
        x: formatTime(value.timestamp)
      });
    }
  }
  if (data.ping.values) {
    for (const value of data.ping.values) {
      pingMin.push({
        y: value.value.minValue,
        x: formatTime(value.timestamp)
      });
      pingMean.push({
        y: value.value.avgValue,
        x: formatTime(value.timestamp)
      });
      pingMax.push({
        y: value.value.maxValue,
        x: formatTime(value.timestamp)
      });
    }
  }
  if (data.jitter.values) {
    for (const value of data.jitter.values) {
      jitter.push({
        y: value.value.value,
        x: formatTime(value.timestamp)
      });
    }
  }

  return {
    responseDetailsDelay,
    responseDetailsResponse,
    pingMin,
    pingMean,
    pingMax,
    jitter,
    name
  };
};

const RouteDetailsComponent: FC<RouteDetailsProps> = ({
  routesChains,
  showDetails
}) => {
  const [value, setValue] = React.useState(0);
  const handleChange = (_event, newValue) => {
    setValue(newValue);
  };
  const data = format(routesChains);
  const ping = [...data.pingMin, ...data.pingMean, ...data.pingMax];
  return (
    <>
      <span
        style={{ display: 'flex', alignItems: 'center' }}
        onClick={() => showDetails(false, routesChains)}
      >
        <ArrowBackIcon color="primary" fontSize="large" />
        <span
          style={{
            fontSize: '1rem',
            fontWeight: 'bold',
            padding: '0 0.4rem',
            display: 'flex',
            alignItems: 'center'
          }}
        >
          {truncate(data.name, 70)}
        </span>
      </span>
      <hr />
      <AppBar position="static">
        <Tabs value={value} onChange={handleChange} indicatorColor="secondary">
          <Tab label="Response length" {...a11yProps(0)} />
          <Tab label="Response delay" {...a11yProps(1)} />
          <Tab label="Ping" {...a11yProps(2)} />
          <Tab label="Jitter" {...a11yProps(3)} />
        </Tabs>
      </AppBar>
      <TabPanel value={value} index={0}>
        <ChartComponent
          name="Response Length"
          values={data.responseDetailsResponse}
          options={RESPONSE_LENGTH_OPTIONS}
        />
      </TabPanel>
      <TabPanel value={value} index={1}>
        <ChartComponent
          name="Response Delay"
          values={data.responseDetailsDelay}
          options={DELAY_OPTIONS}
        />
      </TabPanel>
      <TabPanel value={value} index={2}>
        <ChartComponent name="Ping" values={ping} options={PING_OPTIONS} />
      </TabPanel>
      <TabPanel value={value} index={3}>
        <ChartComponent
          name="Jitter"
          values={data.jitter}
          options={JITTER_OPTIONS}
        />
      </TabPanel>
    </>
  );
};

export default RouteDetailsComponent;
