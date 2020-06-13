import React, { FC, useState, ChangeEvent } from 'react';
import Type from './Groups';
import GridBody from './GridBody';
import { Card, CardContent } from '@material-ui/core';
import InfoOutlinedIcon from '@material-ui/icons/InfoOutlined';
import TextField from '@material-ui/core/TextField';
import FormControlLabel from '@material-ui/core/FormControlLabel';
import Checkbox from '@material-ui/core/Checkbox';
import Favorite from '@material-ui/icons/Favorite';
import FavoriteBorder from '@material-ui/icons/FavoriteBorder';

const requestsTypeSupported = ['get', 'post', 'put', 'delete', 'patch'];

const hyperTexts = ['https', 'http', 'manual'];

const Input: FC<{}> = () => {
  const [requestType, setRequestType] = useState('');
  const [hyperTextType, setHyperTextType] = useState('');
  const [valueURLRoute, setValueURLRoute] = useState('');
  const getRequestType = (type: string) => {
    setRequestType(type);
  };
  const getHyperTextType = (type: string) => {
    console.warn('type is ', type);
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
                color="secondary"
                icon={<FavoriteBorder />}
                checkedIcon={<Favorite />}
              />
            }
            label="Header"
          />
          <FormControlLabel
            control={
              <Checkbox
                color="secondary"
                icon={<FavoriteBorder />}
                checkedIcon={<Favorite />}
              />
            }
            label="Params"
          />
          <FormControlLabel
            control={
              <Checkbox
                color="secondary"
                icon={<FavoriteBorder />}
                checkedIcon={<Favorite />}
              />
            }
            label="Body"
          />
        </div>
        <div style={{ margin: '2%' }}>
          <GridBody name="Header" />
        </div>
      </CardContent>
    </Card>
  );
};

export default Input;
