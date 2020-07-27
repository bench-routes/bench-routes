import React, { useState, ChangeEvent, useEffect } from 'react';
import Type from './Groups';
import GridBody, { pair } from './GridBody';
import { Card, CardContent } from '@material-ui/core';
import InfoOutlinedIcon from '@material-ui/icons/InfoOutlined';
import TextField from '@material-ui/core/TextField';
import FormControlLabel from '@material-ui/core/FormControlLabel';
import Checkbox from '@material-ui/core/Checkbox';
import Button from '@material-ui/core/Button';
import { HOST_IP, paramsTransformValue } from '../../utils/types';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';
import MuiAlert, { AlertProps } from '@material-ui/lab/Alert';
import Snackbar from '@material-ui/core/Snackbar';
import PropTypes from 'prop-types';
import { populateParams } from '../../utils/parse';

interface InputScreenProps {
  screenType: string | undefined;
  params: { Name: string; Value: string }[];
  headers: { OfType: string; Value: string }[];
  route: string;
  body: { Name: string; Value: string }[];
  method: string;
  updateCurrentModal: (routes: any, URL: string) => void;
}

interface AlertSnackBar {
  severity: 'success' | 'error' | 'warning' | 'info';
  message: string;
}

const requestsTypeSupported = ['get', 'post', 'put', 'delete', 'patch'];

const hyperTexts = ['https', 'http', 'manual'];

function Alert(props: AlertProps) {
  return <MuiAlert elevation={6} variant="filled" {...props} />;
}

const Input = (props: InputScreenProps) => {
  const {
    screenType,
    headers,
    params,
    route,
    body,
    method,
    updateCurrentModal
  } = props;
  const [requestType, setRequestType] = useState('');
  const [, setHyperTextType] = useState('');
  const [valueURLRoute, setValueURLRoute] = useState('');

  const [applyHeader, setApplyHeader] = useState<boolean>(false);
  const [headerValues, setHeaderValues] = useState<pair[]>();

  const [applyParams, setApplyParams] = useState<boolean>(false);
  const [paramsValues, setParamsValues] = useState<pair[]>();

  const [applyBody, setApplyBody] = useState<boolean>(false);
  const [bodyValues, setBodyValues] = useState<pair[]>();

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
      setParamsValues(paramValues);
      setBodyValues(bodyValues);
      setHeaderValues(headerValues);
      setRequestType(method);
      setValueURLRoute(route);
    }
  }, [body, headers, method, params, route, screenType]);

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
    setRequestType('');
    setHyperTextType('');
    setValueURLRoute('');
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
    const { route } = props;
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
        body: body
      })
    })
      .then(resp => resp.json())
      .then(response => {
        try {
          const inJSON = JSON.stringify(
            JSON.parse(response.data.ReponseStringified),
            null,
            4
          );
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
            orgRoute: route
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
    setShowResponseButton(false);
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
        body: body
      })
    })
      .then(response => response.json())
      .then(
        response => {
          try {
            const inJSON = JSON.stringify(
              JSON.parse(response['data']),
              null,
              4
            );
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
              body: body
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
                    'error occurred: please contact the dev team or open a issue at https://github.com/zairza-cetb/bench-routes'
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
              'error occurred: please contact the dev team or open a issue at https://github.com/zairza-cetb/bench-routes'
          });
          setShowSnackBar(true);
        }
      );
  };
  return (
    <Card>
      <CardContent>
        <h3 style={{ marginBottom: '2%' }}>Quick Input</h3>
        <h6 style={{ color: 'cadetblue' }}>
          <InfoOutlinedIcon /> Input routes into bench-routes for monitoring
        </h6>
        <hr />
        <div style={{ margin: '2% 0% 2% 0%' }}>
          <h6>HTTP methods</h6>
          <Type slice={requestsTypeSupported} getRequestType={getRequestType} />
        </div>
        <div style={{ margin: '2% 0% 2% 0%' }}>
          <h6>URL</h6>
          <Type slice={hyperTexts} getRequestType={getHyperTextType} />
          <TextField
            id="outlined-basic"
            style={{ margin: '3vh 0vh 0vh 1vh', width: '100%' }}
            value={valueURLRoute}
            onChange={updateURLRouteValue}
            size="medium"
            label="URL route"
            variant="outlined"
          />
        </div>
        <div
          style={{
            border: '1px solid #c4c4c4',
            borderRadius: '1vh',
            padding: '2vh'
          }}
        >
          <h6>Apply</h6>
          <hr />
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
        </div>
        <div style={{ margin: '2%' }}>
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
        <div>
          {screenType === 'config-screen' ? (
            <Button
              variant="contained"
              color="primary"
              onClick={() => testAndEdit()}
            >
              Test & Save
            </Button>
          ) : (
            <Button
              variant="contained"
              color="primary"
              onClick={() => testURL()}
            >
              Test & Save
            </Button>
          )}
          <Button
            variant="contained"
            color="secondary"
            style={{ marginLeft: '1%' }}
            onClick={() => handleCancel()}
          >
            Cancel
          </Button>
          {showResponseButton ? (
            <Button
              variant="contained"
              color="default"
              style={{ marginLeft: '1%' }}
              onClick={() => setOpen(true)}
            >
              Show Response
            </Button>
          ) : null}
        </div>
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
          <Alert severity={openSnackBar.severity}>{openSnackBar.message}</Alert>
        </Snackbar>
      </CardContent>
    </Card>
  );
};

Input.propTypes = {
  screenType: PropTypes.string,
  params: PropTypes.array
};

export default Input;
