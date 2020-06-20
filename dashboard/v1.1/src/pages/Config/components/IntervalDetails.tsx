import React, { useState } from 'react';
import { Grid, TextField, Button, Typography } from '@material-ui/core';
import { HOST_IP } from '../../../utils/types';

const IntervalDetails = (props: any) => {
  const [inputValue, setInputValue] = useState<string>(
    props.durationValue || ''
  );

  const handleIntervalOnChange = e => {
    setInputValue(e.target.value);
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
    <Grid container>
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
        </div>
      ) : (
        <Typography variant="h5">{props.durationValue}</Typography>
      )}
    </Grid>
  );
};

export default IntervalDetails;
