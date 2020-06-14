import React, { FC, useState, ChangeEvent, MouseEvent } from 'react';
import Type from './Groups';
import GridBody, { pair } from './GridBody';
import { Card, CardContent } from '@material-ui/core';
import InfoOutlinedIcon from '@material-ui/icons/InfoOutlined';
import TextField from '@material-ui/core/TextField';
import FormControlLabel from '@material-ui/core/FormControlLabel';
import Checkbox from '@material-ui/core/Checkbox';
import Button from '@material-ui/core/Button';
import URLBuilder from './URLBuilder';

const requestsTypeSupported = ['get', 'post', 'put', 'delete', 'patch'];

const hyperTexts = ['https', 'http', 'manual'];

const Input: FC<{}> = () => {
  const [requestType, setRequestType] = useState('');
  const [hyperTextType, setHyperTextType] = useState('');

  const [valueURLRoute, setValueURLRoute] = useState('');

  const [applyHeader, setApplyHeader] = useState<boolean>(false);
  const [headerValues, setHeaderValues] = useState<pair[]>();

  const [applyParams, setApplyParams] = useState<boolean>(false);
  const [paramsValues, setParamsValues] = useState<pair[]>();

  const [applyBody, setApplyBody] = useState<boolean>(false);
  const [bodyValues, setBodyValues] = useState<pair[]>();

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
  const testURL = () => {
    // const url = new URLBuilder(valueURLRoute);
    if (headerValues === undefined) {
      setHeaderValues([]);
    }
    if (paramsValues === undefined) {
      setParamsValues([]);
    }
    if (bodyValues === undefined) {
      setBodyValues(bodyValues);
    }
    console.warn(headerValues);
    console.warn(paramsValues);
    console.warn(bodyValues);
    const url = new URLBuilder(
      valueURLRoute,
      headerValues,
      paramsValues,
      bodyValues
    );
    url.send(requestType.toLowerCase());
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
                checked={applyHeader}
                onClick={() => setApplyHeader(!applyHeader)}
              />
            }
            label="Header"
          />
          <FormControlLabel
            control={
              <Checkbox
                color="primary"
                checked={applyParams}
                onClick={() => setApplyParams(!applyParams)}
              />
            }
            label="Params"
          />
          <FormControlLabel
            control={
              <Checkbox
                color="primary"
                checked={applyBody}
                onClick={() => setApplyBody(!applyBody)}
              />
            }
            label="Body"
          />
        </div>
        <div style={{ margin: '2%' }}>
          {applyHeader ? (
            <GridBody name="Header" updateParent={setHeaderValues} />
          ) : null}
          {applyParams ? (
            <GridBody name="Params" updateParent={setParamsValues} />
          ) : null}
          {applyBody ? (
            <GridBody name="Body" updateParent={setBodyValues} />
          ) : null}
        </div>
        <div>
          <Button variant="contained" color="primary" onClick={() => testURL()}>
            Test & Save
          </Button>
          <Button
            variant="contained"
            color="secondary"
            style={{ marginLeft: '1%' }}
            onClick={() => {
              setRequestType('');
              setHyperTextType('');
              setValueURLRoute('');
              setHeaderValues([]);
              setParamsValues([]);
              setBodyValues([]);
              setApplyHeader(false);
              setApplyParams(false);
              setApplyBody(false);
            }}
          >
            Cancel
          </Button>
        </div>
      </CardContent>
    </Card>
  );
};

export default Input;
