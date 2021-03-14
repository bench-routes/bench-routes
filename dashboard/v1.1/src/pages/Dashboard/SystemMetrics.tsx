import React, { FC, useEffect, useState } from 'react';
import CPUUsage from './CPUUsage';
import AppBar from '@material-ui/core/AppBar';
import Tabs from '@material-ui/core/Tabs';
import Tab from '@material-ui/core/Tab';
import PropTypes from 'prop-types';
import Box from '@material-ui/core/Box';
import { makeStyles } from '@material-ui/core/styles';
import Alert from '@material-ui/lab/Alert';
import MemoryUsagePercent from './MemoryUsage';
import DiskUsage from './Disk';
import MemoryDetails from './MemoryDetails';
import TimeInstance, { formatTime } from '../../utils/brt';
import { HOST_IP } from '../../utils/types';
import { APIResponse, init } from '../../utils/service';
import { QueryResponse, QueryValues, chartData } from '../../utils/queryTypes';

export function TabPanel(props) {
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

export function a11yProps(index) {
  return {
    id: `simple-tab-${index}`,
    'aria-controls': `simple-tabpanel-${index}`
  };
}

export const useStyles = makeStyles(theme => ({
  root: {
    flexGrow: 1,
    backgroundColor: theme.palette.background.paper
  }
}));

const segregateMetrics = (metricValues: QueryValues[]) => {
  const cpuUsageSlice: chartData[] = [];

  const diskSliceCache: chartData[] = [];
  const diskSliceDiskIO: chartData[] = [];

  const memorySliceAvailableBytes: chartData[] = [];
  const memorySliceFreeBytes: chartData[] = [];
  const memorySliceTotalBytes: chartData[] = [];
  const memorySliceUsedBytes: chartData[] = [];

  const memoryUsedPercentSlice: chartData[] = [];

  for (const metric of metricValues) {
    cpuUsageSlice.push({
      y: metric.value.cpuTotalUsage,
      x: formatTime(metric.timestamp)
    });

    diskSliceCache.push({
      y: metric.value.disk.cached,
      x: formatTime(metric.timestamp)
    });
    diskSliceDiskIO.push({
      y: metric.value.disk.diskIO,
      x: formatTime(metric.timestamp)
    });

    memorySliceAvailableBytes.push({
      y: metric.value.memory.availableBytes,
      x: formatTime(metric.timestamp)
    });
    memorySliceFreeBytes.push({
      y: metric.value.memory.freeBytes,
      x: formatTime(metric.timestamp)
    });
    memorySliceTotalBytes.push({
      y: metric.value.memory.totalBytes,
      x: formatTime(metric.timestamp)
    });
    memorySliceUsedBytes.push({
      y: metric.value.memory.usedBytes,
      x: formatTime(metric.timestamp)
    });

    memoryUsedPercentSlice.push({
      y: metric.value.memory.usedPercent,
      x: formatTime(metric.timestamp)
    });
  }
  return {
    cpuUsageSlice,
    diskSliceCache,
    diskSliceDiskIO,
    memorySliceAvailableBytes,
    memorySliceFreeBytes,
    memorySliceTotalBytes,
    memorySliceUsedBytes,
    memoryUsedPercentSlice
  };
};

interface SystemMetricsProps {
  showLoader(status: boolean): any;
}

const SystemMetrics: FC<SystemMetricsProps> = ({ showLoader }) => {
  const classes = useStyles();
  const [response, setResponse] = useState(init());
  const [error, setError] = useState('');
  const [value, setValue] = React.useState(0);
  const endTimestamp = new Date().getTime() * 1000000 - TimeInstance.Hour;

  useEffect(() => {
    showLoader(true);
  }, [showLoader]);
  useEffect(() => {
    fetch(
      `${HOST_IP}/query?timeSeriesPath=storage/system&endTimestamp=${endTimestamp}`
    )
      .then(res => res.json())
      .then(
        (response: APIResponse<QueryResponse>) => {
          setResponse(response);
        },
        (err: string) => {
          setError(err);
        }
      );
    // eslint-disable-next-line
  }, []);
  const handleChange = (_event, newValue) => {
    setValue(newValue);
  };

  if (error) {
    showLoader(false);
    return <Alert severity="error">Unable to reach the service: error</Alert>;
  }
  if (!response.data.values) {
    return (
      <>
        <Alert severity="info">Fetching data from sources</Alert>
      </>
    );
  }

  const responseInFormat = segregateMetrics(response.data.values);
  showLoader(false);

  return (
    <div className="row">
      <div className="col-md-12" style={{ marginBottom: '1%' }}>
        <div className={classes.root}>
          <AppBar position="static">
            <Tabs
              value={value}
              onChange={handleChange}
              indicatorColor="secondary"
            >
              <Tab
                label="System"
                {...a11yProps(0)}
                style={{ outline: '0px' }}
              />
              <Tab label="Disk" {...a11yProps(1)} style={{ outline: '0px' }} />
              <Tab
                label="Memory details"
                {...a11yProps(2)}
                style={{ outline: '0px' }}
              />
            </Tabs>
          </AppBar>
          <TabPanel value={value} index={0}>
            <div className="row">
              <div className="col-md-6">
                <CPUUsage cpuMetrics={responseInFormat.cpuUsageSlice} />
              </div>
              <div className="col-md-6">
                <MemoryUsagePercent
                  memoryUsagePercentMetrics={
                    responseInFormat.memoryUsedPercentSlice
                  }
                />
              </div>
            </div>
          </TabPanel>
          <TabPanel value={value} index={1}>
            <div className="col-md-12">
              <DiskUsage
                diskIO={responseInFormat.diskSliceDiskIO}
                cache={responseInFormat.diskSliceCache}
              />
            </div>
          </TabPanel>
          <TabPanel value={value} index={2}>
            <div className="col-md-12">
              <MemoryDetails
                availableBytes={responseInFormat.memorySliceAvailableBytes}
                freeBytes={responseInFormat.memorySliceFreeBytes}
                totalBytes={responseInFormat.memorySliceTotalBytes}
                usedBytes={responseInFormat.memorySliceUsedBytes}
              />
            </div>
          </TabPanel>
        </div>
      </div>
    </div>
  );
};

export default SystemMetrics;
