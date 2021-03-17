import React, { useState, ChangeEvent, useEffect } from 'react';
import GridBody, { pair } from './GridBody';
import {
  Card,
  CardContent,
  Container,
  Grid,
  Typography
} from '@material-ui/core';
import MaterialSelect from './components/MaterialSelect';
import MaterialLabelSelector from './components/MaterialLabelSelector';
import TextField from '@material-ui/core/TextField';
import FormControlLabel from '@material-ui/core/FormControlLabel';
import Checkbox from '@material-ui/core/Checkbox';
import Button from '@material-ui/core/Button';
import { HOST_IP, LabelType, paramsTransformValue } from '../../utils/types';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';
import { Alert } from '@material-ui/lab';
import Snackbar from '@material-ui/core/Snackbar';
import PropTypes from 'prop-types';
import { populateLabels, populateParams } from '../../utils/parse';
import {
  defaultHTTPMethodsList,
  defaultRequestMethodsList
} from '../../services/input';
import useStyles from './styles/main';

interface InputScreenProps {
  screenType: string | undefined;
  params: { Name: string; Value: string }[];
  headers: { OfType: string; Value: string }[];
  route: string;
  body: { Name: string; Value: string }[];
  method: string;
  labels?: string[];
  updateCurrentModal: (routes: any, URL: string) => void;
}

interface AlertSnackBar {
  severity: 'success' | 'error' | 'warning' | 'info';
  message: string;
}

const Input = (props: InputScreenProps) => {
  const classes = useStyles();
  const {
    screenType,
    headers,
    params,
    route,
    body,
    method,
    labels,
    updateCurrentModal
  } = props;

  const INITIAL_STATE = {
    REQUEST_TYPE: defaultRequestMethodsList[0].value,
    HTTP_METHOD: defaultHTTPMethodsList[1].value,
    URL_ROUTE: `${defaultHTTPMethodsList[1].value.toLowerCase()}://`
  };

  const [requestType, setRequestType] = useState(INITIAL_STATE.REQUEST_TYPE);
  const [hyperTextType, setHyperTextType] = useState(INITIAL_STATE.HTTP_METHOD);
  const [valueURLRoute, setValueURLRoute] = useState(INITIAL_STATE.URL_ROUTE);

  const [applyHeader, setApplyHeader] = useState<boolean>(false);
  const [headerValues, setHeaderValues] = useState<pair[]>();

  const [applyParams, setApplyParams] = useState<boolean>(false);
  const [paramsValues, setParamsValues] = useState<pair[]>();

  const [applyBody, setApplyBody] = useState<boolean>(false);
  const [bodyValues, setBodyValues] = useState<pair[]>();

  const [selectedLabels, setSelectedLabels] = useState<LabelType[]>([]);

  const [testInputResponse, setTestInputResponse] = useState<string>('');

  const [openSnackBar, setOpenSnackBar] = useState<AlertSnackBar>({
    severity: 'info',
    message: ''
  });
  const [showSnackBar, setShowSnackBar] = useState<boolean>(false);
  const [open, setOpen] = useState(false);
  const [showResponseButton, setShowResponseButton] = useState<boolean>(false);

  useEffect(() => {
    if (screenType === 'config-screen') {
      let paramValues: paramsTransformValue[] = populateParams(params);
      let bodyValues: paramsTransformValue[] = populateParams(body);
      let headerValues: paramsTransformValue[] = populateParams(headers);
      if (labels) {
        const labelValues: LabelType[] = populateLabels(labels);
        setSelectedLabels(labelValues);
      }
      setParamsValues(paramValues);
      setBodyValues(bodyValues);
      setHeaderValues(headerValues);
      setRequestType(method);
      setValueURLRoute(route);
    }
  }, [body, headers, method, params, route, screenType, labels]);

  const getRequestType = (type: string) => {
    setRequestType(type);
  };
  const getHyperTextType = (type: string) => {
    setHyperTextType(type);
    type = type.toLowerCase();
    if (type !== 'manual') {
      setValueURLRoute(`${type}://`);
    } else {
      setValueURLRoute('');
    }
  };
  const updateURLRouteValue = (
    e: ChangeEvent<HTMLTextAreaElement | HTMLInputElement>
  ) => {
    setValueURLRoute(e.target.value);
  };
  const handleCancel = () => {
    setShowResponseButton(false);
    setRequestType(INITIAL_STATE.REQUEST_TYPE);
    setHyperTextType(INITIAL_STATE.HTTP_METHOD);
    setValueURLRoute(INITIAL_STATE.URL_ROUTE);
    setHeaderValues([]);
    setParamsValues([]);
    setBodyValues([]);
    setApplyHeader(false);
    setApplyParams(false);
    setApplyBody(false);
  };
  const testAndEdit = () => {
    const params = {};
    const headers = {};
    const body = {};
    const labels: string[] = [];
    const { route } = props;

    // Is Label selection is mandatory?
    // If yes, then display a dialog
    // to allow user to set them.
    if (selectedLabels.length) {
      selectedLabels.forEach((label: LabelType) => {
        labels.push(label.name);
      });
    }

    setShowResponseButton(false);
    if (headerValues !== undefined) {
      for (const h of headerValues) {
        if (!(h.key === '' && h.value === '')) {
          headers[h.key] = h.value;
        }
      }
    } else {
      setHeaderValues([]);
    }

    if (paramsValues !== undefined) {
      for (const p of paramsValues) {
        if (!(p.key === '' && p.value === '')) {
          params[p.key] = p.value;
        }
      }
    } else {
      setParamsValues([]);
    }

    if (bodyValues !== undefined) {
      for (const b of bodyValues) {
        if (!(b.key === '' && b.value === '')) {
          body[b.key] = b.value;
        }
      }
    } else {
      setBodyValues(bodyValues);
    }
    fetch(`${HOST_IP}/quick-input`, {
      method: 'post',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        method: requestType,
        url: valueURLRoute,
        params: params,
        headers: headers,
        body: body,
        labels
      })
    })
      .then(resp => resp.json())
      .then(response => {
        try {
          const inJSON = JSON.stringify(response['data'], undefined, 4);
          setTestInputResponse(inJSON);
          setShowResponseButton(true);
        } catch (e) {
          setTestInputResponse(response.data.ReponseStringified);
        }
        fetch(`${HOST_IP}/update-route`, {
          method: 'post',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({
            method: requestType,
            url: valueURLRoute,
            params: params,
            headers: headers,
            body: body,
            orgRoute: route,
            labels
          })
        })
          .then(resp => resp.json())
          .then(response => {
            updateCurrentModal(response, valueURLRoute);
            setShowResponseButton(true);
          });
      });
  };
  const testURL = () => {
    const params = {};
    const headers = {};
    const body = {};
    const labels: string[] = [];
    setShowResponseButton(false);

    if (selectedLabels.length) {
      selectedLabels.forEach((label: LabelType) => {
        labels.push(label.name);
      });
    }

    if (headerValues !== undefined) {
      for (const h of headerValues) {
        headers[h.key] = h.value;
      }
    } else {
      setHeaderValues([]);
    }

    if (paramsValues !== undefined) {
      for (const p of paramsValues) {
        params[p.key] = p.value;
      }
    } else {
      setParamsValues([]);
    }

    if (bodyValues !== undefined) {
      for (const b of bodyValues) {
        body[b.key] = b.value;
      }
    } else {
      setBodyValues(bodyValues);
    }

    fetch(`${HOST_IP}/quick-input`, {
      method: 'post',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        method: requestType,
        url: valueURLRoute,
        params: params,
        headers: headers,
        body: body,
        labels
      })
    })
      .then(response => response.json())
      .then(
        response => {
          try {
            const inJSON = JSON.stringify(response['data'], null, 4);
            setTestInputResponse(inJSON);
          } catch (e) {
            setTestInputResponse(response['data']);
          }
          fetch(`${HOST_IP}/add-route`, {
            method: 'post',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
              method: requestType,
              url: valueURLRoute,
              params: params,
              headers: headers,
              body: body,
              labels
            })
          })
            .then(response => response.json())
            .then(
              () => {
                setOpenSnackBar({
                  severity: 'success',
                  message: 'success'
                });
                setShowSnackBar(true);
              },
              err => {
                console.error(err);
                setOpenSnackBar({
                  severity: 'error',
                  message:
                    'error occurred: please contact the dev team or open a issue at https://github.com/bench-routes/bench-routes'
                });
                setShowSnackBar(true);
              }
            );
          setShowResponseButton(true);
        },
        err => {
          console.error(err);
          setOpenSnackBar({
            severity: 'error',
            message:
              'error occurred: please contact the dev team or open a issue at https://github.com/bench-routes/bench-routes'
          });
          setShowSnackBar(true);
        }
      );
  };
  const updateLabels = (labels: LabelType[]) => {
    setSelectedLabels(labels);
  };
  return (
    <>
      <Typography variant="h3" className={classes.pageheader}>
        Quick Input
      </Typography>
      <Container className={classes.container}>
        {/* URL */}
        <Grid container spacing={1} className={classes.controls}>
          <Grid item lg={9} sm={9}>
            <TextField
              id="outlined-basic"
              style={{ width: '100%' }}
              value={valueURLRoute}
              onChange={updateURLRouteValue}
              size="medium"
              label="URL route"
              variant="outlined"
            />
          </Grid>
          <Grid item lg={3} sm={3}>
            {screenType === 'config-screen' ? (
              <Button
                className={classes.btn}
                fullWidth
                variant="contained"
                color="primary"
                onClick={() => testAndEdit()}
              >
                Test and Save
              </Button>
            ) : (
              <Button
                className={classes.btn}
                fullWidth
                variant="contained"
                color="primary"
                onClick={() => testURL()}
              >
                Test and Save
              </Button>
            )}
          </Grid>
        </Grid>
        <Grid
          className={classes.additionalParams}
          container
          spacing={4}
          alignContent="space-between"
        >
          {/* HTTP methods */}
          <Grid item lg={4} sm={4}>
            <MaterialSelect
              variant="outlined"
              dropdownitems={defaultHTTPMethodsList}
              id="http-methods"
              defaultValue={hyperTextType}
              label="HTTP Method"
              onSelection={type => getHyperTextType(type)}
            />
          </Grid>
          {/* Request Type */}
          <Grid item lg={4} sm={4}>
            <MaterialSelect
              variant="outlined"
              dropdownitems={defaultRequestMethodsList}
              id="http-methods"
              defaultValue={requestType}
              label="Request Type"
              onSelection={type => getRequestType(type)}
            />
          </Grid>
          {/* Labels */}
          <Grid item lg={4} sm={4}>
            <MaterialLabelSelector
              defaultLabels={selectedLabels}
              updateLabels={labels => updateLabels(labels)}
            />
          </Grid>
        </Grid>
        <Grid className={classes.additionalParams}>
          <Button
            className={classes.marginRtSm}
            variant="contained"
            color="primary"
            disableElevation
            onClick={() => handleCancel()}
          >
            Reset
          </Button>
          {showResponseButton ? (
            <Button
              variant="contained"
              color="primary"
              disableElevation
              onClick={() => setOpen(true)}
            >
              Show Response
            </Button>
          ) : null}
          <Dialog aria-labelledby="customized-dialog-title" open={open}>
            <DialogTitle id="customized-dialog-title">Response</DialogTitle>
            <DialogContent dividers>
              <Card>
                <CardContent>
                  <pre style={{ fontWeight: 'bold' }}>{testInputResponse}</pre>
                </CardContent>
              </Card>
            </DialogContent>
            <DialogActions>
              <Button onClick={() => setOpen(false)} color="secondary">
                Close
              </Button>
            </DialogActions>
          </Dialog>
          <Snackbar
            open={showSnackBar}
            autoHideDuration={6000}
            onClose={() => setShowSnackBar(false)}
          >
            <Alert
              elevation={6}
              variant="filled"
              severity={openSnackBar.severity}
            >
              {openSnackBar.message}
            </Alert>
          </Snackbar>
        </Grid>
      </Container>
      {/* Headers, params and others */}
      <Container className={classes.container}>
        <Grid container className={classes.marginTopMd}>
          <Grid item lg={12} sm={12}>
            <FormControlLabel
              control={
                <Checkbox
                  color="primary"
                  checked={applyHeader || (headerValues || []).length > 0}
                  onClick={() => setApplyHeader(!applyHeader)}
                />
              }
              label="Header"
            />
            <FormControlLabel
              control={
                <Checkbox
                  color="primary"
                  checked={applyParams || (paramsValues || []).length > 0}
                  onClick={() => setApplyParams(!applyParams)}
                />
              }
              label="Params"
            />
            <FormControlLabel
              control={
                <Checkbox
                  color="primary"
                  checked={applyBody || (bodyValues || []).length > 0}
                  onClick={() => setApplyBody(!applyBody)}
                />
              }
              label="Body"
            />
            <div className={classes.params}>
              {applyHeader || (headerValues || []).length > 0 ? (
                <GridBody
                  name="Header"
                  headers={headerValues}
                  updateParent={setHeaderValues}
                />
              ) : null}
              {applyParams || (paramsValues || []).length > 0 ? (
                <GridBody
                  name="Params"
                  headers={paramsValues}
                  updateParent={setParamsValues}
                />
              ) : null}
              {applyBody || (bodyValues || []).length > 0 ? (
                <GridBody
                  name="Body"
                  headers={bodyValues}
                  updateParent={setBodyValues}
                />
              ) : null}
            </div>
          </Grid>
        </Grid>
      </Container>
    </>
  );
};

Input.propTypes = {
  screenType: PropTypes.string,
  params: PropTypes.array
};

export default Input;
