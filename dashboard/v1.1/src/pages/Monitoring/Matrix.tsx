import React, { FC, useState } from 'react';
import { useFetch } from '../../utils/useFetch';
import { HOST_IP } from '../../utils/types';
import { QueryResponse } from '../../utils/queryTypes';
import { columns } from './Columns';

import TableContainer from '@material-ui/core/TableContainer';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import CircularProgress from '@material-ui/core/CircularProgress';

import { Badge } from 'reactstrap';

interface Path {
  fping: string;
  jitter: string;
  monitor: string;
  ping: string;
  matrixName: string;
}

export interface TimeSeriesPath {
  name: string;
  path: Path;
}

interface MatrixProps {
  timeSeriesPath: TimeSeriesPath[];
}

interface ElementProps {
  timeSeriesPath: TimeSeriesPath;
}

interface MatrixResponse {
  jitter: QueryResponse;
  monitor: QueryResponse;
  ping: QueryResponse;
}

const round = (n: string): number => {
  const num = parseInt(n, 10);
  return Math.round(num * 10) / 10;
};

const Element: FC<ElementProps> = ({ timeSeriesPath }) => {
  const { response, error } = useFetch<MatrixResponse>(
    `${HOST_IP}/query-matrix?routeNameMatrix=${timeSeriesPath.path.matrixName}`
  );
  const [, render] = useState();
  if (error) {
    return null;
  }

  if (!response.data) {
    return (
      <TableRow className="table-data-row" key={timeSeriesPath.path.matrixName}>
        <TableCell style={{ minWidth: 170, fontSize: 16 }} align="left">
          <Badge color="light">{timeSeriesPath.name}</Badge>
        </TableCell>
        <TableCell
          style={{ minWidth: 100, fontSize: 16 }}
          align="center"
        ></TableCell>
        <TableCell
          style={{ minWidth: 170, fontSize: 16 }}
          align="center"
        ></TableCell>
        <TableCell style={{ minWidth: 170 }} align="center"></TableCell>
        <TableCell style={{ minWidth: 170 }} align="center"></TableCell>
        <TableCell style={{ minWidth: 170, fontSize: 16 }} align="center">
          <CircularProgress size={20} thickness={6} />
        </TableCell>
      </TableRow>
    );
  }

  console.warn(response.data);
  const data = response.data;
  // setTimeout(() => {
  //   render({});
  // }, 10 * 1000);
  return (
    <TableRow>
      <TableCell style={{ minWidth: 170, fontSize: 16 }} align="left">
        <Badge color="light">{timeSeriesPath.name}</Badge>
      </TableCell>
      <TableCell style={{ minWidth: 100, fontSize: 16 }} align="center">
        <Badge color="warning">
          {round(data.ping.values[0].value.avgValue)} ms
        </Badge>
      </TableCell>
      <TableCell style={{ minWidth: 170, fontSize: 16 }} align="center">
        <Badge color="warning">
          {round(data.jitter.values[0].value.value)} ms
        </Badge>
      </TableCell>
      <TableCell style={{ minWidth: 170, fontSize: 16 }} align="center">
        <Badge color="warning">{data.monitor.values[0].value.delay} ms</Badge>
      </TableCell>
      <TableCell style={{ minWidth: 170, fontSize: 16 }} align="center">
        <Badge color="warning">{data.monitor.values[0].value.resLength}</Badge>
      </TableCell>
      <TableCell style={{ minWidth: 170, fontSize: 16 }} align="center">
        <Badge color="success">{'UP'}</Badge>
      </TableCell>
    </TableRow>
  );
};

const Matrix: FC<MatrixProps> = ({ timeSeriesPath }) => (
  <TableContainer>
    <Table stickyHeader>
      <TableHead>
        <TableRow>
          {columns.map((column, i) => (
            <TableCell
              key={i}
              align={column.align}
              style={{
                minWidth: column.minWidth,
                fontWeight: 600,
                fontSize: 15
              }}
            >
              {column.label}
            </TableCell>
          ))}
        </TableRow>
      </TableHead>
      <TableBody>
        {timeSeriesPath.map((series, i) => (
          <Element timeSeriesPath={series} key={i} />
        ))}
      </TableBody>
    </Table>
  </TableContainer>
);

export default Matrix;
