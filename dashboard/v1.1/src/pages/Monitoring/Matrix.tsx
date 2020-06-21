import React, { FC, useState, useEffect } from 'react';
import { HOST_IP } from '../../utils/types';
import {
  TimeSeriesPath,
  MatrixResponse,
  RouteDetails
} from '../../utils/queryTypes';
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
import Tooltip from '@material-ui/core/Tooltip';

import { Badge } from 'reactstrap';

type APIResponse<MatrixResponse> = { status: string; data: MatrixResponse };

interface MatrixProps {
  timeSeriesPath: TimeSeriesPath[];
  showRouteDetails(status: boolean, details: RouteDetails): void;
}

interface ElementProps {
  timeSeriesPath: TimeSeriesPath;
  showRouteDetails(status: boolean, details: RouteDetails): void;
}

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
      async function fetchDetails() {
        try {
          const response = await fetch(
            `${HOST_IP}/query-matrix?routeNameMatrix=${instance.path.matrixName}&endTimestamp=${endTimestamp}`
          );
          const matrix = (await response.json()) as APIResponse<RouteDetails>;
          resolve(matrix.data);
        } catch (e) {
          console.error(e);
          showWarning(true);
          reject(e);
        }
      }
      fetchDetails();
    });

    const details = await monitoringDetails;
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
        <TableCell
          style={{ width: 170, fontSize: 16, overflowX: 'hidden' }}
          align="left"
        >
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
  }, 10 * 1000);

  return (
    <TableRow>
      <TableCell
        style={{ maxWidth: 240, fontSize: 16, overflowX: 'hidden' }}
        align="left"
      >
        <Badge color="light">
          <Tooltip title={timeSeriesPath.name}>
            <div>{timeSeriesPath.name}</div>
          </Tooltip>
        </Badge>
      </TableCell>
      <TableCell style={{ minWidth: 100, fontSize: 16 }} align="center">
        <Badge color="warning">
          {data.ping === undefined ? (
            '-'
          ) : (
            <>
              {data.ping.values === null
                ? '-'
                : round(data.ping.values[0].value.avgValue)}
            </>
          )}{' '}
          ms
        </Badge>
      </TableCell>
      <TableCell style={{ minWidth: 120, fontSize: 16 }} align="center">
        <Badge color="warning">
          {data.jitter === undefined ? (
            '-'
          ) : (
            <>
              {data.jitter.values === null
                ? '-'
                : round(data.jitter.values[0].value.value)}
            </>
          )}{' '}
          ms
        </Badge>
      </TableCell>
      <TableCell style={{ minWidth: 150, fontSize: 16 }} align="center">
        <Badge color="warning">
          {data.monitor === undefined ? (
            '-'
          ) : (
            <>
              {' '}
              {data.monitor.values === null
                ? '-'
                : data.monitor.values[0].value.delay}{' '}
            </>
          )}{' '}
          ms
        </Badge>
      </TableCell>
      <TableCell style={{ minWidth: 170, fontSize: 16 }} align="center">
        <Badge color="warning">
          {data.monitor === undefined ? (
            '-'
          ) : (
            <>
              {data.monitor.values === null
                ? '-'
                : data.monitor.values[0].value.resLength}
            </>
          )}
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
