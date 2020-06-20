import React, { FC, useState, useEffect } from 'react';
import { HOST_IP } from '../../utils/types';
import IntervalDetails from './components/IntervalDetails';
import {
  Grid,
  Card,
  CardContent,
  Typography,
  makeStyles,
  Chip,
  Tooltip
} from '@material-ui/core';
import { Edit as EditIcon, Close as CloseIcon } from '@material-ui/icons';
import SearchTable from './components/MaterialTable';
import { truncate } from '../../utils/stringManipulations';

interface IntervalType {
  test: string;
  duration: number;
  unit: string;
}

interface TableRouteType {
  route: string;
  methods: string[];
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
  const [configRoutes, setConfigRoutes] = useState<[string, string[]][] | null>(
    null
  );
  const [toggleResults, setToggleResults] = useState({
    ping: false,
    jitter: false,
    'req-res-delay-and-monitoring': false
  });

  const fetchConfigIntervals = async () => {
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
  };

  const fetchConfigRoutes = async () => {
    const response = await fetch(`${HOST_IP}/get-config-routes`).then(resp => {
      return resp.json();
    });
    const { data } = response;
    let configRoutes = new Map();
    data.forEach(route => {
      const uri = route['URL'] + '/' + route['Route'];
      if (!configRoutes.has(uri)) {
        configRoutes.set(uri, [route['Method']]);
      } else {
        configRoutes.set(uri, [...configRoutes.get(uri), route['Method']]);
      }
    });
    setConfigRoutes(Array.from(configRoutes));
  };

  useEffect(() => {
    fetchConfigIntervals().then(() => fetchConfigRoutes());
  }, []);

  const getTableRoutes = (routes: [string, string[]][] | null) => {
    let tableData: TableRouteType[] = [];
    routes?.forEach(r => {
      tableData.push({
        route: r[0],
        methods: r[1]
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
      render: (rowData: { methods: any[]; route: string }) =>
        rowData.methods.map(m => (
          <Chip
            key={rowData.route + m}
            variant="outlined"
            color="primary"
            label={m}
          />
        ))
    }
  ];

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
      <SearchTable
        title=""
        columns={columns}
        data={getTableRoutes(configRoutes)}
        editable={{
          onRowUpdate: (newData, oldData) =>
            new Promise((resolve, reject) => {
              setTimeout(() => {
                // const index = oldData.tableData.id;

                resolve();
              }, 1000);
            }),
          onRowDelete: oldData =>
            new Promise((resolve, reject) => {
              setTimeout(() => {
                // const index = oldData.tableData.id;

                resolve();
              }, 1000);
            })
        }}
      />
    </>
  );
};

export default Config;
