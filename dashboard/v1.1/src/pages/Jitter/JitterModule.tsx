import React, { FC, useState, useEffect } from 'react';
import { useFetch } from '../../utils/useFetch';
import { service_states, HOST_IP } from '../../utils/types';
import { Card, CardContent, TextField } from '@material-ui/core';
import {
  RoutesSummary,
  QueryResponse,
  chartData,
  APIQueryResponse,
  APITimeSeriesResponse,
  PathData,
  TimeSeries
} from '../../utils/queryTypes';
import Autocomplete from '@material-ui/lab/Autocomplete';
import Alert from '@material-ui/lab/Alert';
import { formatTime } from '../../utils/brt';
import Jitter from '../Monitoring/Jitter';

const format = (datas: QueryResponse) => {
  const jitter: chartData[] = [];

  for (const value of datas.values) {
    jitter.push({
      y: value.value.value,
      x: formatTime(value.timestamp)
    });
  }

  return {
    jitter
  };
};

interface showChartsDataParam {
  jitter: chartData[];
}

const JitterModule: FC<{}> = () => {
  const [routesDetails, setRoutesDetails] = useState<RoutesSummary>();
  const [value] = useState(routesDetails?.testServicesRoutes[0]);
  const [inputValue, setInputValue] = useState('');
  const [showCharts, setShowCharts] = useState(false);
  const [jitterData, setJitterData] = useState<showChartsDataParam>();

  const { response, error } = useFetch<service_states>(
    `${HOST_IP}/service-state`
  );

  useEffect(() => {
    fetch(`${HOST_IP}/routes-summary`)
      .then(res => res.json())
      .then((response: { status: string; data: RoutesSummary }) => {
        setRoutesDetails(response.data);
      });
  }, []);

  async function getChartsData(v: string) {
    try {
      const response = await fetch(`${HOST_IP}/get-route-time-series`);
      const TimeSeriesArray = (await response.json()) as APITimeSeriesResponse;
      const SpecificRoute = TimeSeriesArray.data.find(
        (item: { name: string; path: PathData }) => item.name === v
      ) as TimeSeries;

      const response2 = await fetch(
        `${HOST_IP}/query?timeSeriesPath=${SpecificRoute.path.jitter}`
      );
      const matrix = (await response2.json()) as APIQueryResponse;
      var formatdata = format(matrix.data);
      setJitterData(formatdata);
      setShowCharts(true);
    } catch (e) {
      console.error(e);
    }
  }

  if (error) {
    return <Alert severity="error">Unable to reach the service: error</Alert>;
  }
  if (!response.data) {
    return <Alert severity="info">Fetching from sources</Alert>;
  }

  // const states: service_states = response.data;
  const options =
    routesDetails?.testServicesRoutes !== undefined
      ? routesDetails.testServicesRoutes
      : ['Please fill routes'];

  return (
    <Card>
      <CardContent>
        <div>
          <h4>Jitter</h4>
          <div style={{ float: 'right', marginTop: '-45px' }}>
            <Autocomplete
              value={value}
              onChange={(event, newValue) => {}}
              inputValue={inputValue}
              onInputChange={(event, newInputValue) => {
                setInputValue(newInputValue);
                getChartsData(newInputValue);
              }}
              id="controllable-states-demo"
              options={options}
              style={{ width: 300 }}
              renderInput={params => (
                <TextField
                  {...params}
                  label="Select Route"
                  variant="outlined"
                />
              )}
            />
          </div>
        </div>
        <br />
        <hr />
        <div>
          {jitterData !== undefined && showCharts ? (
            <Jitter value={jitterData.jitter} />
          ) : null}
        </div>
      </CardContent>
    </Card>
  );
};

export default JitterModule;
