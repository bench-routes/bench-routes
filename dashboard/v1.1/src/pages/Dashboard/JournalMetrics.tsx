import React, { FC, useEffect, useState, useContext } from 'react';
import Chart from 'react-apexcharts';
import Alert from '@material-ui/lab/Alert';
import { Tabs, Tab, AppBar } from '@material-ui/core';
import { APIResponse, init } from '../../utils/service';
import { QueryResponse, QueryValues, chartData } from '../../utils/queryTypes';
import { HOST_IP } from '../../utils/types';
import { formatTime } from '../../utils/brt';
import { useStyles, TabPanel, a11yProps } from './SystemMetrics';
import TimePanel from './TimePanel';
import { ThemeContext } from '../../layouts/BaseLayout';

const format = (data: QueryValues[] | any) => {
  const cerr: chartData[] = [];
  const cwarn: chartData[] = [];
  const cevents: chartData[] = [];
  const ckerr: chartData[] = [];
  const ckwarn: chartData[] = [];
  const ckevents: chartData[] = [];
  if (!data) {
    return {
      cerr,
      cwarn,
      cevents,
      ckerr,
      ckwarn,
      ckevents
    };
  }
  for (const value of data) {
    cerr.push({
      y: value.value ? value.value.cerr : null,
      x: formatTime(value.timestamp)
    });
    cwarn.push({
      y: value.value ? value.value.cwarn : null,
      x: formatTime(value.timestamp)
    });
    cevents.push({
      y: value.value ? value.value.cevents : null,
      x: formatTime(value.timestamp)
    });
    ckerr.push({
      y: value.value ? value.value.ckerr : null,
      x: formatTime(value.timestamp)
    });
    ckwarn.push({
      y: value.value ? value.value.ckwarn : null,
      x: formatTime(value.timestamp)
    });
    ckevents.push({
      y: value.value ? value.value.ckevents : null,
      x: formatTime(value.timestamp)
    });
  }
  return {
    cerr,
    cwarn,
    cevents,
    ckerr,
    ckwarn,
    ckevents
  };
};

interface JournalMetricsProps {
  startTimestamp?: number;
  endTimestamp?: number;
}

const JournalMetrics: FC<JournalMetricsProps> = ({
  startTimestamp,
  endTimestamp
}) => {
  const classes = useStyles();
  const themeMode = useContext(ThemeContext);
  const [response, setResponse] = useState(init());
  const [fetchTime, setfetchTime] = useState(0);
  const [error, setError] = useState('');
  const [value, setValue] = React.useState(0);
  const handleChange = (_event, newValue) => {
    setValue(newValue);
  };

  const fetchDetails = async (): Promise<QueryResponse> => {
    return new Promise<QueryResponse>(async (resolve, reject) => {
      try {
        const response = await fetch(
          endTimestamp
            ? `${HOST_IP}/query?timeSeriesPath=storage/journal&startTimestamp=${endTimestamp}&endTimestamp=${startTimestamp}`
            : `${HOST_IP}/query?timeSeriesPath=storage/journal&endTimestamp=${startTimestamp}`
        );
        const data = (await response.json()) as APIResponse<QueryResponse>;
        resolve(data.data);
      } catch (error) {
        reject(error);
      }
    });
  };

  const fetchSystemData = async () => {
    try {
      const start = performance.now();
      const details = await fetchDetails();
      const end = performance.now();
      setResponse(details);
      setfetchTime(end - start);
    } catch (error) {
      setError(error);
    }
  };

  useEffect(() => {
    setResponse(init());
    fetchSystemData();
    // eslint-disable-next-line
  }, [endTimestamp, startTimestamp]);
  if (error) {
    return <Alert severity="error">Unable to reach the service: error</Alert>;
  }
  if (!response.values.length) {
    return <Alert severity="info">Fetching data from sources</Alert>;
  }
  const data = format(response.values);
  const seriesSystemd = [
    {
      name: 'Errors',
      data: data.cerr
    },
    {
      name: 'Warnings',
      data: data.cwarn
    },
    {
      name: 'Events',
      data: data.cevents
    }
  ];
  const seriesKernel = [
    {
      name: 'Errors',
      data: data.ckerr
    },
    {
      name: 'Warnings',
      data: data.ckwarn
    },
    {
      name: 'Events',
      data: data.ckevents
    }
  ];
  const optionsSystemd = {
    chart: {
      type: 'area'
    },
    dataLabels: {
      enabled: false
    },
    stroke: {
      show: true,
      curve: 'straight',
      lineCap: 'butt',
      width: 1
    },
    subtitle: {
      text: 'Systemd services',
      align: 'center'
    },
    fill: {
      opacity: 1,
      type: 'gradient',
      gradient: {
        shade: 'dark',
        type: 'vertical',
        shadeIntensity: 0.3,
        inverseColors: true,
        opacityFrom: 0.8,
        opacityTo: 0.2
      }
    },
    theme: {
      mode: themeMode
    }
  };
  const optionsKernel = {
    chart: {
      type: 'area'
    },
    dataLabels: {
      enabled: false
    },
    stroke: {
      show: true,
      curve: 'straight',
      lineCap: 'butt',
      width: 1
    },
    subtitle: {
      text: 'Kernel',
      align: 'center'
    },
    fill: {
      opacity: 1,
      type: 'gradient',
      gradient: {
        shade: 'dark',
        type: 'vertical',
        shadeIntensity: 0.3,
        inverseColors: true,
        opacityFrom: 0.8,
        opacityTo: 0.2
      }
    },
    theme: {
      mode: themeMode
    }
  };

  return (
    <div className={classes.root}>
      <div
        style={{
          display: 'flex',
          padding: '0 0 16px 0'
        }}
      >
        <TimePanel
          fetchTime={fetchTime}
          evaluationTime={response.evaluationTime}
        />
      </div>
      <AppBar position="static">
        <Tabs value={value} onChange={handleChange} indicatorColor="secondary">
          <Tab label="Kernel" {...a11yProps(0)} />
          <Tab label="Systemd" {...a11yProps(1)} />
        </Tabs>
      </AppBar>
      <TabPanel value={value} index={0}>
        <div className="row">
          <div className="col-md-12">
            <Chart
              series={seriesKernel}
              options={optionsKernel}
              height="300"
              type="area"
            />
          </div>
        </div>
      </TabPanel>
      <TabPanel value={value} index={1}>
        <div className="row">
          <div className="col-md-12">
            <Chart
              series={seriesSystemd}
              options={optionsSystemd}
              height="300"
              type="area"
            />
          </div>
        </div>
      </TabPanel>
    </div>
  );
};

export default JournalMetrics;
