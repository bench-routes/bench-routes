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
import ChartComponent from '../Monitoring/Chart';
import { JITTER_OPTIONS } from '../../utils/constants/chart';

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
    let res = filterUrl(v);
    res = res.substring(0, res.indexOf('/'));

    try {
      const response = await fetch(
        `${HOST_IP}/query?timeSeriesPath=storage/jitter/chunk_jitter_${res}`
      );
      const matrix = (await response.json()) as APIQueryResponse;
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
            <ChartComponent
              name="Jitter"
              values={jitterData.jitter}
              options={JITTER_OPTIONS}
            />
          ) : null}
        </div>
      </CardContent>
    </Card>
  );
};

export default JitterModule;
