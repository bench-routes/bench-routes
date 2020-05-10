import React, { FC } from 'react';
import { useFetch } from '../../utils/useFetch';
import CPUUsage from './CPUUsage';
import ExpansionPanel from '@material-ui/core/ExpansionPanel';
import ExpansionPanelDetails from '@material-ui/core/ExpansionPanelDetails';
import ExpansionPanelSummary from '@material-ui/core/ExpansionPanelSummary';
import ExpandMoreIcon from '@material-ui/icons/ExpandMore';
import Typography from '@material-ui/core/Typography';
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

const useStyles = makeStyles(theme => ({
  heading: {
    fontSize: theme.typography.pxToRem(15),
    flexBasis: '33.33%',
    flexShrink: 0,
    fontWeight: 600
  },
  secondaryHeading: {
    fontSize: theme.typography.pxToRem(15),
    color: theme.palette.text.secondary
  }
}));

interface SystemMetricsProps {
  done(status: boolean): any;
}

const SystemMetrics: FC<SystemMetricsProps> = ({ done }) => {
  const classes = useStyles();
  const endTimestamp = new Date().getTime() * 1000000 - TimeInstance.Hour;
  const { response, error } = useFetch<QueryResponse>(
    `${HOST_IP}/query?timeSeriesPath=storage/system&endTimestamp=${endTimestamp}`
  );
  if (error) {
    console.warn(error);
  }
  if (!response.data) {
    return null;
  }
  done(true);
  const responseInFormat = segregateMetrics(response.data.values);
  return (
    <div className="row">
      <div className="col-md-12" style={{ marginBottom: '1%' }}>
        <ExpansionPanel>
          <ExpansionPanelSummary
            expandIcon={<ExpandMoreIcon />}
            aria-controls="panel1bh-content"
            id="panel1bh-header"
          >
            <Typography className={classes.heading}>
              System performance
            </Typography>
            <Typography className={classes.secondaryHeading}>
              Performance values related to central processing
            </Typography>
          </ExpansionPanelSummary>
          <ExpansionPanelDetails>
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
          </ExpansionPanelDetails>
        </ExpansionPanel>
      </div>
      <div className="col-md-12" style={{ marginBottom: '1%' }}>
        <ExpansionPanel>
          <ExpansionPanelSummary
            expandIcon={<ExpandMoreIcon />}
            aria-controls="panel1bh-content"
            id="panel1bh-header"
          >
            <Typography className={classes.heading}>
              Disk performance
            </Typography>
            <Typography className={classes.secondaryHeading}>
              Performance values related to system-disk
            </Typography>
          </ExpansionPanelSummary>
          <ExpansionPanelDetails>
            <div className="col-md-12">
              <DiskUsage metrics={responseInFormat.diskSlice} />
            </div>
          </ExpansionPanelDetails>
        </ExpansionPanel>
      </div>
      <div className="col-md-12" style={{ marginBottom: '1%' }}>
        <ExpansionPanel>
          <ExpansionPanelSummary
            expandIcon={<ExpandMoreIcon />}
            aria-controls="panel1bh-content"
            id="panel1bh-header"
          >
            <Typography className={classes.heading}>Memory details</Typography>
            <Typography className={classes.secondaryHeading}>
              Detail visualization of memory values
            </Typography>
          </ExpansionPanelSummary>
          <ExpansionPanelDetails>
            <div className="col-md-12">
              <MemoryDetails metrics={responseInFormat.memorySlice} />
            </div>
          </ExpansionPanelDetails>
        </ExpansionPanel>
      </div>
    </div>
  );
};

export default SystemMetrics;
