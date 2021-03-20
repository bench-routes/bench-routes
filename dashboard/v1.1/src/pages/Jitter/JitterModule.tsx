import React, { FC, useState, useEffect } from 'react';
import { useFetch } from '../../utils/useFetch';
import { service_states, HOST_IP } from '../../utils/types';
import { Card, CardContent, TextField } from '@material-ui/core';
import {
  QueryResponse,
  chartData,
  APIQueryResponse,
  APITimeSeriesResponse,
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
  const [hashRoutMap, SetHashRouteMap] = useState(new Map());
  const [routesDetails, setRoutesDetails] = useState<string[]>();
  const [value] = useState(routesDetails !== undefined ? routesDetails : '');
  const [inputValue, setInputValue] = useState('');
  const [showCharts, setShowCharts] = useState(false);
  const [renderError, setRenderError] = useState(false);
  const [jitterData, setJitterData] = useState<showChartsDataParam>();

  const { response, error } = useFetch<service_states>(
    `${HOST_IP}/service-state`
  );

  useEffect(() => {
    fetch(`${HOST_IP}/get-route-time-series`)
      .then(res => res.json())
      .then((response: APITimeSeriesResponse) => {
        let tempMap: Map<string, string> = new Map();
        let tempArr: string[] = [];

        response.data.forEach((item: TimeSeries) => {
          tempArr.push(item.name);
          tempMap.set(item.name, item.path.matrixName);
        });

        SetHashRouteMap(tempMap);
        setRoutesDetails(tempArr);
      });
  }, []);

  async function getChartsData(v: string) {
    try {
      const response = await fetch(
        `${HOST_IP}/query?timeSeriesPath=storage/jitter/chunk_jitter_${hashRoutMap.get(
          v
        )}`
      );
      const matrix = (await response.json()) as APIQueryResponse;
      var formatdata = format(matrix.data);
      setJitterData(formatdata);
      setShowCharts(true);
      setRenderError(false);
    } catch (e) {
      console.error(e);
      setRenderError(true);
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
    routesDetails !== undefined ? routesDetails : ['Please fill routes'];

  return (
    <Card>
      <CardContent>
        <div>
          <h4>Jitter</h4>
          <div style={{ float: 'right', marginTop: '-45px' }}>
            <Autocomplete
              value={value[0]}
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
          {!renderError ? (
            jitterData !== undefined && showCharts ? (
              <Jitter value={jitterData.jitter} />
            ) : null
          ) : (
            <Alert severity="error">No data found</Alert>
          )}
        </div>
      </CardContent>
    </Card>
  );
};

export default JitterModule;
