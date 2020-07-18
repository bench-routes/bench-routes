import React, { FC, useState, useEffect, lazy, Suspense } from 'react';
import {
  HOST_IP,
  service_states,
  routeEntryType,
  routeOptionsInterface
} from '../../utils/types';
import IntervalDetails from './components/IntervalDetails';
import {
  Grid,
  Card,
  CardContent,
  Typography,
  makeStyles,
  Chip,
  Tooltip,
  CircularProgress
} from '@material-ui/core';
import { Edit as EditIcon, Close as CloseIcon } from '@material-ui/icons';
import { truncate } from '../../utils/stringManipulations';
import { useFetch } from '../../utils/useFetch';
import { Alert } from '@material-ui/lab';
import { tableIcons } from '../../utils/tableIcons';
import EditModal from './components/EditModal';

const SearchTable = lazy(() => import('./components/MaterialTable'));

interface IntervalType {
  test: string;
  duration: number;
  unit: string;
}

interface TableRouteType {
  route: string;
  methods: string[];
}

interface TableRowData {
  methods: string[];
  route: string;
  tableData: { id: number };
}

const useStyles = makeStyles(theme => ({
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

const Config: FC<{}> = () => {
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
    'req-res-delay-and-monitoring': false
  });
  const [editModalOpen, setEditModalOpen] = useState(false);
  const [selectedRow, setSelectedRow] = useState<routeEntryType>({
    route: '',
    options: []
  });

  const { response, error } = useFetch<service_states>(
    `${HOST_IP}/service-state`
  );

  const fetchConfigIntervals = async () => {
    try {
      const response = await fetch(`${HOST_IP}/get-config-intervals`).then(
        resp => {
          return resp.json();
        }
      );
      const { data } = response;
      let intervals: any = [];
      data.forEach(interval => {
        intervals.push({
          test: interval['Test'],
          duration: interval['Duration'],
          unit: interval['Type']
        });
      });
      setConfigIntervals(intervals);
    } catch (e) {
      console.log(e);
    }
  };

  const fetchConfigRoutes = async () => {
    const response = await fetch(`${HOST_IP}/get-config-routes`).then(resp => {
      return resp.json();
    });
    const { data } = response;
    let configRoutes = new Map();
    data.forEach(route => {
      const uri = route['URL'];
      if (!configRoutes.has(uri)) {
        configRoutes.set(uri, [
          {
            method: route['Method'],
            body: route['Body'],
            headers: route['Header'],
            params: route['Params']
          }
        ]);
      } else {
        configRoutes.set(uri, [
          ...configRoutes.get(uri),
          {
            method: route['Method'],
            body: route['Body'],
            headers: route['Header'],
            params: route['Params']
          }
        ]);
      }
    });
    setConfigRoutes(configRoutes);
  };

  useEffect(() => {
    fetchConfigIntervals().then(() => fetchConfigRoutes());
  }, []);

  const getTableData = (routes: [string, routeOptionsInterface[]][]) => {
    let tableData: TableRouteType[] = [];
    routes.forEach(route => {
      let methods: string[] = [];
      route[1].forEach(option => {
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

  const columns = [
    {
      field: 'route',
      title: 'Route'
    },
    {
      field: 'methods',
      title: 'Methods',
      render: (rowData: TableRowData) =>
        rowData.methods.map(m => (
          <Chip
            key={rowData.route + m}
            variant="outlined"
            color="primary"
            label={m}
            style={{ marginLeft: '2px' }}
          />
        ))
    }
  ];

  const updateConfigRoutes = routes => {
    const { data } = routes;
    let configRoutes = new Map();
    data.forEach(route => {
      const uri = route['URL'];
      if (!configRoutes.has(uri)) {
        configRoutes.set(uri, [
          {
            method: route['Method'],
            body: route['Body'],
            headers: route['Header'],
            params: route['Params']
          }
        ]);
      } else {
        configRoutes.set(uri, [
          ...configRoutes.get(uri),
          {
            method: route['Method'],
            body: route['Body'],
            headers: route['Header'],
            params: route['Params']
          }
        ]);
      }
    });
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
      case 'req-res-delay-and-monitoring':
        setToggleResults({
          ...toggleResults,
          'req-res-delay-and-monitoring': !toggleResults[
            'req-res-delay-and-monitoring'
          ]
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
      <Grid container spacing={4}>
        {configIntervals?.map(interval => {
          const { test, duration, unit } = interval;
          return (
            <Grid item lg={3} sm={6} xl={3} xs={12} key={test}>
              <Card className={classes.cardStyle}>
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
                    reFetch={fetchConfigIntervals}
                    toggleComponentView={(name: string) => handleToggle(name)}
                    toggleResults={toggleResults}
                    durationValue={duration}
                    intervalName={test}
                  />
                  <Typography variant="body1">{unit}</Typography>
                </CardContent>
              </Card>
            </Grid>
          );
        })}
      </Grid>
      <EditModal
        isOpen={editModalOpen}
        setOpen={(open: boolean) => setEditModalOpen(open)}
        selectedRoute={selectedRow}
        updateConfigRoutes={route => updateConfigRoutes(route)}
      />
      <Suspense fallback={<CircularProgress disableShrink />}>
        <SearchTable
          title=""
          columns={columns}
          data={getTableData(Array.from(configRoutes))}
          actions={[
            {
              icon: tableIcons.Edit,
              tooltip: 'Edit Route',
              onClick: (event, rowData: TableRowData) => {
                console.log(configRoutes.get(rowData.route));
                setSelectedRow({
                  route: rowData.route,
                  options: configRoutes.get(rowData.route)
                });
                setEditModalOpen(!editModalOpen);
              }
            }
          ]}
          editable={{
            onRowDelete: async (oldData: TableRouteType) => {
              try {
                await fetch(`${HOST_IP}/delete-route`, {
                  method: 'post',
                  headers: { 'Content-Type': 'application/json' },
                  body: JSON.stringify({
                    actualRoute: oldData.route
                  })
                })
                  .then(resp => resp.json())
                  .then(response => {
                    const { data } = response;
                    let configRoutes = new Map();
                    data.forEach(route => {
                      const uri = route['URL'];
                      if (!configRoutes.has(uri)) {
                        configRoutes.set(uri, [
                          {
                            method: route['Method'],
                            body: route['Body'],
                            headers: route['Header'],
                            params: route['Params']
                          }
                        ]);
                      } else {
                        configRoutes.set(uri, [
                          ...configRoutes.get(uri),
                          {
                            method: route['Method'],
                            body: route['Body'],
                            headers: route['Header'],
                            params: route['Params']
                          }
                        ]);
                      }
                    });
                    setConfigRoutes(configRoutes);
                  });
              } catch (e) {
                console.log(e);
              }
            }
          }}
        />
      </Suspense>
    </>
  );
};

export default Config;
