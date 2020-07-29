import React, { useState, useEffect, lazy, Suspense } from 'react';
import {
  HOST_IP,
  service_states,
  routeEntryType,
  routeOptionsInterface
} from '../../utils/types';
import IntervalDetails from './components/IntervalDetails';
import {
  Grid,
  CardContent,
  Typography,
  makeStyles,
  Tooltip,
  CircularProgress
} from '@material-ui/core';
import Paper from '@material-ui/core/Paper';
import { Edit as EditIcon, Close as CloseIcon } from '@material-ui/icons';
import { truncate } from '../../utils/stringManipulations';
import { useFetch } from '../../utils/useFetch';
import { Alert } from '@material-ui/lab';
import { tableIcons } from '../../utils/tableIcons';
import EditModal from './components/EditModal';
import { getRoutesMap } from '../../utils/parse';
import { fetchConfigIntervals, fetchConfigRoutes } from '../../services/config';
import { columns, TableRowData } from './components/MaterialTable';
import { onRowDelete, TableRouteType, IntervalType } from './handles';

const SearchTable = lazy(() => import('./components/MaterialTable'));

const useStyles = makeStyles(_theme => ({
  root: {
    diplay: 'flex'
  },
  cardStyle: {
    minHeight: '8vh'
  },
  h6: {
    fontWeight: 'normal'
  }
}));

const Config = () => {
  const classes = useStyles();
  const [configIntervals, setConfigIntervals] = useState<IntervalType[] | null>(
    null
  );

  const [configRoutes, setConfigRoutes] = useState<
    Map<string, routeOptionsInterface[]>
  >(new Map());

  const [toggleResults, setToggleResults] = useState({
    ping: false,
    jitter: false,
    monitoring: false
  });

  const [editModalOpen, setEditModalOpen] = useState(false);

  const [selectedRow, setSelectedRow] = useState<routeEntryType>({
    route: '',
    options: []
  });

  const { response, error } = useFetch<service_states>(
    `${HOST_IP}/service-state`
  );

  useEffect(() => {
    fetchConfigIntervals(setConfigIntervals).then(() =>
      fetchConfigRoutes(setConfigRoutes)
    );
  }, []);

  const getTableData = (routes: [string, routeOptionsInterface[]][]) => {
    let tableData: TableRouteType[] = [];
    routes.forEach((route: [string, routeOptionsInterface[]]) => {
      let methods: string[] = [];
      route[1].forEach((option: routeOptionsInterface) => {
        Object.keys(option).forEach(k => {
          if (k === 'method') {
            methods.push(option[k]);
          }
        });
      });
      tableData.push({
        route: route[0],
        methods
      });
    });
    return tableData;
  };

  const updateConfigRoutes = routes => {
    const { data } = routes;
    let configRoutes: Map<string, routeOptionsInterface[]> = getRoutesMap(data);
    setConfigRoutes(configRoutes);
  };

  const handleToggle = (intervalName: string) => {
    switch (intervalName) {
      case 'ping':
        setToggleResults({ ...toggleResults, ping: !toggleResults['ping'] });
        break;
      case 'jitter':
        setToggleResults({
          ...toggleResults,
          jitter: !toggleResults['jitter']
        });
        break;
      case 'monitoring':
        setToggleResults({
          ...toggleResults,
          monitoring: !toggleResults['monitoring']
        });
        break;
    }
  };

  if (error) {
    return <Alert severity="error">Unable to reach the service: error</Alert>;
  }
  if (!response.data) {
    return <Alert severity="info">Fetching from sources</Alert>;
  }

  return (
    <>
      <Paper elevation={0} style={{ marginBottom: '1%', padding: '1%' }}>
        <CardContent>
          <h4>Scrape Intervals</h4>
          <hr />
        </CardContent>
        <Grid container spacing={1}>
          {configIntervals?.map((interval: IntervalType) => {
            const { test, duration, unit } = interval;
            return (
              <Grid item lg={3} sm={6} xl={3} xs={12} key={test}>
                <Paper
                  className={classes.cardStyle}
                  elevation={0}
                  variant="outlined"
                >
                  <CardContent>
                    <Grid container style={{ justifyContent: 'space-between' }}>
                      <Grid item>
                        <Typography
                          gutterBottom
                          variant="h6"
                          className={classes.h6}
                        >
                          {truncate(
                            test.charAt(0).toUpperCase() + test.slice(1),
                            14
                          )}
                        </Typography>
                      </Grid>
                      <Grid item>
                        {toggleResults[test] ? (
                          <Tooltip title="Cancel" style={{ cursor: 'pointer' }}>
                            <CloseIcon onClick={() => handleToggle(test)} />
                          </Tooltip>
                        ) : (
                          <Tooltip title="Edit" style={{ cursor: 'pointer' }}>
                            <EditIcon onClick={() => handleToggle(test)} />
                          </Tooltip>
                        )}
                      </Grid>
                    </Grid>
                    <IntervalDetails
                      reFetch={() => fetchConfigIntervals(setConfigIntervals)}
                      toggleComponentView={(name: string) => handleToggle(name)}
                      toggleResults={toggleResults}
                      durationValue={duration}
                      intervalName={test}
                    />
                    <Typography variant="body1" align="center">
                      {unit}
                    </Typography>
                  </CardContent>
                </Paper>
              </Grid>
            );
          })}
        </Grid>
      </Paper>
      <EditModal
        isOpen={editModalOpen}
        setOpen={(open: boolean) => setEditModalOpen(open)}
        selectedRoute={selectedRow}
        updateConfigRoutes={route => updateConfigRoutes(route)}
      />
      <Suspense fallback={<CircularProgress disableShrink />}>
        <Paper elevation={0}>
          <SearchTable
            title="URL endpoints"
            columns={columns}
            data={getTableData(Array.from(configRoutes))}
            actions={[
              {
                icon: tableIcons.Edit,
                tooltip: 'Edit Route',
                onClick: (_event: any, rowData: TableRowData) => {
                  setSelectedRow({
                    route: rowData.route,
                    options: configRoutes.get(rowData.route)
                  });
                  setEditModalOpen(!editModalOpen);
                }
              }
            ]}
            editable={{ onRowDelete }}
          />
        </Paper>
      </Suspense>
    </>
  );
};

export default Config;
