import React, { FC, useState, useEffect } from 'react';
import { useFetch } from '../../utils/useFetch';
import { service_states, HOST_IP } from '../../utils/types';
import { Card, CardContent, TextField } from '@material-ui/core';
import {
  RoutesSummary,
  QueryResponse,
  chartData,
  APIQueryResponse
} from '../../utils/queryTypes';
import Autocomplete from '@material-ui/lab/Autocomplete';
import Alert from '@material-ui/lab/Alert';
import { formatTime } from '../../utils/brt';
import { filterUrl } from '../../utils/filterUrl';
import { PING_OPTIONS } from '../../utils/constants/chart';
import ChartComponent from '../Monitoring/Chart';

const format = (datas: QueryResponse) => {
  const pingMin: chartData[] = [];
  const pingMean: chartData[] = [];
  const pingMax: chartData[] = [];

  for (const value of datas.values) {
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

  const ping = [...pingMin, ...pingMean, ...pingMax];
  return ping;
};

const PingModule: FC<{}> = () => {
  const [routesDetails, setRoutesDetails] = useState<RoutesSummary>();
  const [value] = useState(routesDetails?.testServicesRoutes[0]);
  const [inputValue, setInputValue] = useState('');
  const [showCharts, setShowCharts] = useState(false);
  const [pingData, setPingData] = useState<chartData[]>();

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
    let res = filterUrl(v);
    res = res.substring(0, res.indexOf('/'));

    try {
      const response = await fetch(
        `${HOST_IP}/query?timeSeriesPath=storage/ping/chunk_ping_${res}`
      );
      const matrix = (await response.json()) as APIQueryResponse;
      var formatdata = format(matrix.data);
      setPingData(formatdata);
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

  // TODO: add the status icon for the module
  // const states: service_states = response.data;

  const options =
    routesDetails?.testServicesRoutes !== undefined
      ? routesDetails.testServicesRoutes
      : ['Please fill routes'];

  return (
    <Card>
      <CardContent>
        <div>
          <h4>Ping</h4>
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
          {pingData !== undefined && showCharts ? (
            <ChartComponent
              name="Ping"
              values={pingData}
              options={PING_OPTIONS}
            />
          ) : null}
        </div>
      </CardContent>
    </Card>
  );
};

export default PingModule;
