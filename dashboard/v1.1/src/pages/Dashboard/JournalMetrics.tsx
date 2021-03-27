import React, { FC, useEffect, useState, useContext } from 'react';
import Chart from 'react-apexcharts';
import Alert from '@material-ui/lab/Alert';
import AppBar from '@material-ui/core/AppBar';
import Tabs from '@material-ui/core/Tabs';
import Tab from '@material-ui/core/Tab';
import { APIResponse, init } from '../../utils/service';
import { QueryResponse, QueryValues, chartData } from '../../utils/queryTypes';
import { HOST_IP } from '../../utils/types';
import { formatTime } from '../../utils/brt';
import { useStyles, TabPanel, a11yProps } from './SystemMetrics';
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
      y: value.value.cerr,
      x: formatTime(value.timestamp)
    });
    cwarn.push({
      y: value.value.cwarn,
      x: formatTime(value.timestamp)
    });
    cevents.push({
      y: value.value.cevents,
      x: formatTime(value.timestamp)
    });
    ckerr.push({
      y: value.value.ckerr,
      x: formatTime(value.timestamp)
    });
    ckwarn.push({
      y: value.value.ckwarn,
      x: formatTime(value.timestamp)
    });
    ckevents.push({
      y: value.value.ckevents,
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
  const [error, setError] = useState('');
  const [value, setValue] = React.useState(0);
  const handleChange = (_event, newValue) => {
    setValue(newValue);
  };
  useEffect(() => {
    setResponse(init());
    fetch(
      endTimestamp
        ? `${HOST_IP}/query?timeSeriesPath=storage/journal&startTimestamp=${endTimestamp}&endTimestamp=${startTimestamp}`
        : `${HOST_IP}/query?timeSeriesPath=storage/journal&endTimestamp=${startTimestamp}`
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
  }, [endTimestamp, startTimestamp]);
  if (error) {
    return <Alert severity="error">Unable to reach the service: error</Alert>;
  }
  if (!response.data.values) {
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
