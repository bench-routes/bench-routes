import React, { FC } from 'react';
import { useFetch } from '../../utils/useFetch';
import CPUUsage from './CPUUsage';
import AppBar from '@material-ui/core/AppBar';
import Tabs from '@material-ui/core/Tabs';
import Tab from '@material-ui/core/Tab';
import PropTypes from 'prop-types';
import Box from '@material-ui/core/Box';
import { makeStyles } from '@material-ui/core/styles';
import MemoryUsagePercent from './MemoryUsage';
import DiskUsage from './Disk';
import MemoryDetails from './MemoryDetails';
import TimeInstance from '../../utils/brt';
import { HOST_IP } from '../../utils/types';
import {
  QueryResponse,
  QueryValues,
  queryValueCPUUsage,
  queryValueMemory,
  queryValueDisk,
  queryValueMemoryUsedPercent
} from '../../utils/queryTypes';

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

function a11yProps(index) {
  return {
    id: `simple-tab-${index}`,
    'aria-controls': `simple-tabpanel-${index}`
  };
}

const useStyles = makeStyles(theme => ({
  root: {
    flexGrow: 1,
    backgroundColor: theme.palette.background.paper
  }
}));

const segregateMetrics = (metricValues: QueryValues[]) => {
  const cpuUsageSlice: queryValueCPUUsage[] = [];
  const diskSlice: queryValueDisk[] = [];
  const memorySlice: queryValueMemory[] = [];
  const memoryUsedPercentSlice: queryValueMemoryUsedPercent[] = [];

  for (const metric of metricValues) {
    cpuUsageSlice.push({
      CPUUsage: metric.value.cpuTotalUsage,
      normalizedTime: metric.normalizedTime
    });
    diskSlice.push({
      cached: metric.value.disk.cached,
      diskIO: metric.value.disk.diskIO,
      normalizedTime: metric.normalizedTime
    });
    memorySlice.push({
      availableBytes: metric.value.memory.availableBytes,
      freeBytes: metric.value.memory.freeBytes,
      totalBytes: metric.value.memory.totalBytes,
      usedBytes: metric.value.memory.usedBytes,
      usedPercent: metric.value.memory.usedPercent,
      normalizedTime: metric.normalizedTime
    });
    memoryUsedPercentSlice.push({
      memoryUsedPercent: metric.value.memory.usedPercent,
      normalizedTime: metric.normalizedTime
    });
  }
  return { cpuUsageSlice, diskSlice, memorySlice, memoryUsedPercentSlice };
};

interface SystemMetricsProps {
  showLoader(status: boolean): any;
}

const SystemMetrics: FC<SystemMetricsProps> = ({ showLoader }) => {
  const classes = useStyles();
  const [value, setValue] = React.useState(0);

  const handleChange = (event, newValue) => {
    setValue(newValue);
  };

  const endTimestamp = new Date().getTime() * 1000000 - TimeInstance.Hour;
  showLoader(true);
  const { response, error } = useFetch<QueryResponse>(
    `${HOST_IP}/query?timeSeriesPath=storage/system&endTimestamp=${endTimestamp}`
  );
  if (error) {
    console.warn(error);
  }
  if (!response.data) {
    return null;
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
              <Tab label="System" {...a11yProps(0)} />
              <Tab label="Disk" {...a11yProps(1)} />
              <Tab label="Memory details" {...a11yProps(2)} />
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
              <DiskUsage metrics={responseInFormat.diskSlice} />
            </div>
          </TabPanel>
          <TabPanel value={value} index={2}>
            <div className="col-md-12">
              <MemoryDetails metrics={responseInFormat.memorySlice} />
            </div>
          </TabPanel>
        </div>
      </div>
    </div>
  );
};

export default SystemMetrics;
