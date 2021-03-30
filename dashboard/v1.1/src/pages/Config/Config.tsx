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
import { handleRowDelete, TableRouteType, IntervalType } from './handles';
import InputLabel from '@material-ui/core/InputLabel';
import MenuItem from '@material-ui/core/MenuItem';
import FormControl from '@material-ui/core/FormControl';
import Select from '@material-ui/core/Select';
import { useXticks } from '../../utils/useXticks';

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
  },
  formControl: {
    margin: 0,
    minWidth: 150
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

  const [tableData, setTableData] = useState<TableRouteType[]>([]);

  const xticks = useXticks();

  useEffect(() => {
    fetchConfigIntervals(setConfigIntervals)
      .then(async () => {
        const routes = await fetchConfigRoutes(setConfigRoutes);
        return routes;
      })
      .then(routes => {
        getTableData(Array.from(routes));
      });
  }, []);

  const handleXticksChange = (event: React.ChangeEvent<{ value: any }>) => {
    const changeXticks = xticks['updateXticks'];
    changeXticks(event.target.value as string);
  };

  const xtickVal: any = [
    '5',
    '10',
    '15',
    '20',
    '25',
    '30',
    '35',
    '40',
    '45',
    '50',
    '55',
    '60',
    '65',
    '70'
  ];

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
    setTableData(tableData);
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
      <Paper elevation={0} style={{ marginBottom: '1%', padding: '1%' }}>
        <CardContent>
          <h4>Chart Options</h4>
          <hr />
        </CardContent>
        <CardContent>
          <FormControl variant="outlined" className={classes.formControl}>
            <InputLabel id="demo-simple-select-outlined-label">
              x-tick amount
            </InputLabel>
            <Select
              labelId="demo-simple-select-outlined-label"
              id="demo-simple-select-outlined"
              value={xticks['xticks']}
              onChange={handleXticksChange}
              label="x-tick amount"
            >
              {xtickVal.map(val => (
                <MenuItem key={val} value={val}>
                  {val}
                </MenuItem>
              ))}
            </Select>
          </FormControl>
        </CardContent>
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
            data={tableData}
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
            editable={{
              onRowDelete: oldData =>
                handleRowDelete(
                  oldData,
                  setConfigRoutes,
                  tableData,
                  setTableData
                )
            }}
          />
        </Paper>
      </Suspense>
    </>
  );
};

export default Config;
