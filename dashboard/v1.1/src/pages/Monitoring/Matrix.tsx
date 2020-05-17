import React, { FC, useState, useEffect } from 'react';
import { HOST_IP } from '../../utils/types';
import { QueryResponse } from '../../utils/queryTypes';
import { columns } from './Columns';
import TimeInstance from '../../utils/brt';

import TableContainer from '@material-ui/core/TableContainer';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import CircularProgress from '@material-ui/core/CircularProgress';
import WarningOutlinedIcon from '@material-ui/icons/WarningOutlined';
import ArrowForwardIcon from '@material-ui/icons/ArrowForward';
import Slide from '@material-ui/core/Slide';

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
  showRouteDetails(status: boolean, details: RouteDetails): void;
}

interface ElementProps {
  timeSeriesPath: TimeSeriesPath;
  showRouteDetails(status: boolean, details: RouteDetails): void;
}

interface MatrixResponse {
  jitter: QueryResponse;
  monitor: QueryResponse;
  ping: QueryResponse;
}

export interface RouteDetails {
  ping: Promise<QueryResponse | null>;
  jitter: Promise<QueryResponse | null>;
  monitor: Promise<QueryResponse | null>;
}

type APIResponse<MatrixResponse> = { status: string; data: MatrixResponse };

const round = (n: string): number => {
  const num = parseInt(n, 10);
  return Math.round(num * 10) / 10;
};

const Pad: FC<{}> = () => <span>&nbsp;&nbsp;&nbsp;&nbsp;</span>;

const Element: FC<ElementProps> = ({ timeSeriesPath, showRouteDetails }) => {
  const [data, setData] = useState<MatrixResponse>();
  const [trigger, setTrigger] = useState<number>(0);
  const [updating, setUpdating] = useState<boolean>(true);
  const [warning, showWarning] = useState<boolean>(false);
  const endTimestamp = new Date().getTime() * 1000000 - TimeInstance.Hour;

  const fetchTimeSeriesDetails = async (instance: TimeSeriesPath) => {
    const monitoringDetails = new Promise<RouteDetails>((resolve, reject) => {
      async function fetchDetails(url: string) {
        try {
          const response = await fetch(url);
          return ((await response.json()) as APIResponse<QueryResponse>).data;
        } catch (e) {
          console.error(e);
          showWarning(true);
          reject(e);
          return null;
        }
      }

      const batch: RouteDetails = {
        ping: fetchDetails(
          `${HOST_IP}/query?timeSeriesPath=${instance.path.ping}&endTimestamp=${endTimestamp}`
        ),
        jitter: fetchDetails(
          `${HOST_IP}/query?timeSeriesPath=${instance.path.jitter}&endTimestamp=${endTimestamp}`
        ),
        monitor: fetchDetails(
          `${HOST_IP}/query?timeSeriesPath=${instance.path.monitor}&endTimestamp=${endTimestamp}`
        )
      };
      resolve(batch);
    });

    const details = await monitoringDetails;
    console.warn('details are', details);
    showRouteDetails(true, details);
  };

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
  }, [trigger, timeSeriesPath.path.matrixName]);

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
            <Pad />
          )}
        </TableCell>
      </TableRow>
    );
  }

  setTimeout(() => {
    setTrigger(trigger + 1);
  }, 15 * 1000);

  return (
    <Slide direction="right" in={true} mountOnEnter timeout={1000}>
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
            // sizes are kept in accordance with the ArrowForwardIcon. Do not change them.
            <CircularProgress disableShrink size={19} thickness={4} />
          ) : (
            <span onClick={() => fetchTimeSeriesDetails(timeSeriesPath)}>
              <ArrowForwardIcon color="primary" />
            </span>
          )}
        </TableCell>
      </TableRow>
    </Slide>
  );
};

const Matrix: FC<MatrixProps> = ({ timeSeriesPath, showRouteDetails }) => (
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
          <Element
            timeSeriesPath={series}
            showRouteDetails={showRouteDetails}
            key={i}
          />
        ))}
      </TableBody>
    </Table>
  </TableContainer>
);

export default Matrix;
