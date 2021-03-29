import React, { FC, useEffect, useState } from 'react';
import Chart from 'react-apexcharts';
import Alert from '@material-ui/lab/Alert';
import AppBar from '@material-ui/core/AppBar';
import Tabs from '@material-ui/core/Tabs';
import Tab from '@material-ui/core/Tab';
import { APIResponse, init } from '../../utils/service';
import { QueryResponse, QueryValues, chartData } from '../../utils/queryTypes';
import { HOST_IP } from '../../utils/types';
import TimeInstance, { formatTime } from '../../utils/brt';
import { useStyles, TabPanel, a11yProps } from './SystemMetrics';

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
    console.log(value.value)
    cerr.push({
      y: (value.value ? value.value.cerr : null),
      x: formatTime(value.timestamp)
    });
    cwarn.push({
      y: (value.value ? value.value.cwarn : null),
      x: formatTime(value.timestamp)
    });
    cevents.push({
      y: (value.value ? value.value.cevents : null),
      x: formatTime(value.timestamp)
    });
    ckerr.push({
      y: (value.value ? value.value.ckerr : null),
      x: formatTime(value.timestamp)
    });
    ckwarn.push({
      y: (value.value ? value.value.ckwarn : null),
      x: formatTime(value.timestamp)
    });
    ckevents.push({
      y: (value.value ? value.value.ckevents : null),
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
  darkMode(status: boolean): any;
}

const JournalMetrics: FC<JournalMetricsProps> = ({ darkMode }) => {
  const classes = useStyles();
  const [response, setResponse] = useState(init());
  const [error, setError] = useState('');
  const endTimestamp = new Date().getTime() * 1000000 - TimeInstance.Hour;
  const [value, setValue] = React.useState(0);
  const handleChange = (_event, newValue) => {
    setValue(newValue);
  };
  useEffect(() => {
    fetch(
      `${HOST_IP}/query?timeSeriesPath=storage/journal&endTimestamp=${endTimestamp}`
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
  if (error) {
    return <Alert severity="error">Unable to reach the service: error</Alert>;
  }
  if (!response.data) {
    return <Alert severity="info">Fetching data from sources</Alert>;
  }
  const data = format(response.data.values);

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
      type: 'area',
      background: '#fff'
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
    tooltip: {
      theme: !darkMode ? 'light' : 'dark'
    }
  };
  const optionsKernel = {
    chart: {
      type: 'area',
      background: '#fff'
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
    tooltip: {
      theme: !darkMode ? 'light' : 'dark'
    }
  };
  return (
    <div className={classes.root}>
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
