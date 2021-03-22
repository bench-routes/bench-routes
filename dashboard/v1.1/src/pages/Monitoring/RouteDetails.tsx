import React, { FC } from 'react';
import {
  chartData,
  EvaluationTime,
  RouteDetails
} from '../../utils/queryTypes';
import { formatTime } from '../../utils/brt';
import ResLength from './ResLength';
import Delay from './Delay';
import Ping from './Ping';
import Jitter from './Jitter';
import AppBar from '@material-ui/core/AppBar';
import Tabs from '@material-ui/core/Tabs';
import Tab from '@material-ui/core/Tab';
import PropTypes from 'prop-types';
import Box from '@material-ui/core/Box';
import ArrowBackIcon from '@material-ui/icons/ArrowBack';
import { Badge } from 'reactstrap';

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

const TabPanel = props => {
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
};

TabPanel.propTypes = {
  children: PropTypes.node,
  index: PropTypes.any.isRequired,
  value: PropTypes.any.isRequired
};

const TimePanel = ({ data, value }) => {
  return (
    <div style={{ borderLeft: '3px solid #2196f3' }}>
      <h6
        style={{
          margin: 0,
          padding: 10,
          display: 'flex',
          alignItems: 'center'
        }}
      >
        Query fetched in
        <Badge color="success" style={{ fontSize: '0.97rem', marginLeft: 10 }}>
          {data.fetchTime.toFixed(3)}ms
        </Badge>
      </h6>
      {value && data.evaluationTime[Type(value)] !== '' ? (
        <h6
          style={{
            margin: 0,
            padding: 10,
            display: 'flex',
            alignItems: 'center'
          }}
        >
          Query Time :
          <Badge
            color="success"
            style={{ fontSize: '0.97rem', marginLeft: 10 }}
          >
            {data.evaluationTime[Type(value)]}
          </Badge>
        </h6>
      ) : null}
    </div>
  );
};

const Type = (val: number) => {
  switch (val) {
    case 1:
      return 'monitor';
    case 2:
      return 'ping';
    case 3:
      return 'jitter';
    default:
      return '';
  }
};

const format = (data: RouteDetails) => {
  const responseDetailsDelay: chartData[] = [];
  const responseDetailsResponse: chartData[] = [];
  const pingMin: chartData[] = [];
  const pingMean: chartData[] = [];
  const pingMax: chartData[] = [];
  const jitter: chartData[] = [];
  const name = data.name;
  const fetchTime = data.fetchTime;
  const evaluationTime: EvaluationTime = { ping: '', jitter: '', monitor: '' };
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
    evaluationTime.monitor = data.monitor.evaluationTime;
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
    evaluationTime.ping = data.ping.evaluationTime;
  }
  if (data.jitter.values) {
    for (const value of data.jitter.values) {
      jitter.push({
        y: value.value.value,
        x: formatTime(value.timestamp)
      });
    }
    evaluationTime.jitter = data.jitter.evaluationTime;
  }

  return {
    responseDetailsDelay,
    responseDetailsResponse,
    pingMin,
    pingMean,
    pingMax,
    jitter,
    name,
    fetchTime,
    evaluationTime
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
  return (
    <>
      <div
        style={{
          display: 'flex',
          justifyContent: 'space-between',
          alignItems: 'center',
          position: 'relative'
        }}
      >
        <span style={{ display: 'flex', position: 'relative', width: '70%' }}>
          <ArrowBackIcon
            color="primary"
            fontSize="large"
            style={{ cursor: 'pointer ' }}
            onClick={() => showDetails(false, routesChains)}
          />
          <h3
            style={{
              fontSize: '1rem',
              fontWeight: 'bold',
              padding: '0.4rem 0.4rem 0 0.4rem',
              margin: 0,
              maxWidth: '90%',
              overflow: 'hidden',
              textOverflow: 'ellipsis',
              whiteSpace: 'nowrap'
            }}
          >
            {data.name}
          </h3>
        </span>
        <TimePanel data={data} value={value} />
      </div>
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
        <ResLength resLength={data.responseDetailsResponse} />
      </TabPanel>
      <TabPanel value={value} index={1}>
        <Delay delay={data.responseDetailsDelay} />
      </TabPanel>
      <TabPanel value={value} index={2}>
        <Ping min={data.pingMin} mean={data.pingMean} max={data.pingMax} />
      </TabPanel>
      <TabPanel value={value} index={3}>
        <Jitter value={data.jitter} />
      </TabPanel>
    </>
  );
};

export default RouteDetailsComponent;
