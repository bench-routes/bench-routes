import React, { useState } from 'react';
import { TextField, Button, Typography } from '@material-ui/core';
import Alert from '@material-ui/lab/Alert';
import { HOST_IP } from '../../../utils/types';

const IntervalDetails = (props: any) => {
  const [inputValue, setInputValue] = useState<Number>(
    props.durationValue || ''
  );

  const [alert, setAlert] = useState<boolean>(false);

  const handleIntervalOnChange = e => {
    const value = e.target.value;
    if (value >= 0) {
      setInputValue(value);
      setAlert(false);
    } else {
      setAlert(true);
    }
  };

  const handleSubmit = async (e, intervalName: string) => {
    await fetch(`${HOST_IP}/config/update-interval`, {
      method: 'post',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        intervalName,
        value: inputValue
      })
    })
      .then(res => {
        return res.json();
      })
      .then(response => {
        if (response.status === '200') {
          props.toggleComponentView(intervalName);
          props.reFetch();
        }
      });
  };
  return (
    <>
      {props.toggleResults[props.intervalName] ? (
        <div>
          <form
            onSubmit={e => handleSubmit(e, props.intervalName)}
            style={{ display: 'flex' }}
          >
            <TextField
              id="outlined-basic"
              label={props.durationValue}
              variant="outlined"
              onChange={e => handleIntervalOnChange(e)}
            />
            <Button
              variant="contained"
              color="primary"
              style={{ marginLeft: '4px' }}
              type="submit"
            >
              Go
            </Button>
          </form>
          <div>
            {alert && (
              <Alert severity="error">Please try a positive int Value</Alert>
            )}
          </div>
        </div>
      ) : (
        <div style={{ textAlign: 'center' }}>
          <Typography variant="h2" align="center">
            {props.durationValue}
          </Typography>
        </div>
      )}
    </>
  );
};

export default IntervalDetails;
