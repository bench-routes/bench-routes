import React, { FC, useState, useEffect } from 'react';
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
import WarningOutlinedIcon from '@material-ui/icons/WarningOutlined';

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

type APIResponse<MatrixResponse> = { status: string; data?: MatrixResponse };

const round = (n: string): number => {
  const num = parseInt(n, 10);
  return Math.round(num * 10) / 10;
};

const Pad: FC<{}> = () => <span>&nbsp;&nbsp;&nbsp;&nbsp;</span>;

const Element: FC<ElementProps> = ({ timeSeriesPath }) => {
  const [data, setData] = useState<MatrixResponse>();
  const [trigger, setTrigger] = useState<number>(0);
  const [updating, setUpdating] = useState<boolean>(true);
  const [warning, showWarning] = useState<boolean>(false);

  useEffect(() => {
    async function fetchMatrix(name: string) {
      try {
        setUpdating(true);
        const response = await fetch(
          `${HOST_IP}/query-matrix?routeNameMatrix=${name}`
        );
        const inMatrixResponse = (await response.json()) as APIResponse<
          MatrixResponse
        >;
        setData(inMatrixResponse.data);
        setTimeout(() => {
          setUpdating(false);
          showWarning(false);
        }, 1000);
      } catch (e) {
        console.error(e);
        showWarning(true);
        return null;
      }
    }
    fetchMatrix(timeSeriesPath.path.matrixName);
  }, [trigger]);

  if (!data) {
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
        <TableCell
          style={{ minWidth: 170, fontSize: 16 }}
          align="center"
        ></TableCell>
        <TableCell style={{ minWidth: 10, fontSize: 16 }} align="center">
          {updating ? (
            <CircularProgress disableShrink size={15} thickness={4} />
          ) : (
            '&nbsp'
          )}
        </TableCell>
      </TableRow>
    );
  }

  setTimeout(() => {
    setTrigger(trigger + 1);
  }, 5 * 1000);

  return (
    <TableRow>
      <TableCell style={{ minWidth: 170, fontSize: 16 }} align="left">
        <Badge color="light">{timeSeriesPath.name}</Badge>
      </TableCell>
      <TableCell style={{ minWidth: 100, fontSize: 16 }} align="center">
        <Badge color="warning">
          {data.ping === undefined
            ? '-'
            : round(data.ping.values[0].value.avgValue)}{' '}
          ms
        </Badge>
      </TableCell>
      <TableCell style={{ minWidth: 170, fontSize: 16 }} align="center">
        <Badge color="warning">
          {data.jitter === undefined
            ? '-'
            : round(data.jitter.values[0].value.value)}{' '}
          ms
        </Badge>
      </TableCell>
      <TableCell style={{ minWidth: 170, fontSize: 16 }} align="center">
        <Badge color="warning">
          {data.monitor === undefined
            ? '-'
            : data.monitor.values[0].value.delay}{' '}
          ms
        </Badge>
      </TableCell>
      <TableCell style={{ minWidth: 170, fontSize: 16 }} align="center">
        <Badge color="warning">
          {data.monitor === undefined
            ? '-'
            : data.monitor.values[0].value.resLength}
        </Badge>
      </TableCell>
      <TableCell style={{ minWidth: 170, fontSize: 16 }} align="center">
        <Badge color="success">{'UP'}</Badge>
      </TableCell>
      <TableCell style={{ minWidth: 10, fontSize: 16 }} align="center">
        {warning ? (
          <WarningOutlinedIcon />
        ) : updating ? (
          <CircularProgress disableShrink size={15} thickness={4} />
        ) : (
          <Pad />
        )}
      </TableCell>
    </TableRow>
  );
};

const Matrix: FC<MatrixProps> = ({ timeSeriesPath }) => (
  <TableContainer style={{ maxHeight: '100vh', overflowY: 'hidden' }}>
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
