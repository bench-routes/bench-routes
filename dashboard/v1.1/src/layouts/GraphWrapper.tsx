import { makeStyles } from '@material-ui/core';
import ArrowBackIcon from '@material-ui/icons/ArrowBack';
import Alert from '@material-ui/lab/Alert';
import moment from 'moment';
import React, { FC, useState } from 'react';
import Datetime from 'react-datetime';
import { truncate } from '../utils/stringManipulations';

interface WrapperProps {
  child: React.ReactElement;
  isMonitoring: boolean;
  routesChains?: { name: string; route: string };
  showDetails?: (
    status: boolean,
    details: { name: string; route: string }
  ) => void;
}

const useStyles = makeStyles(theme => ({
  inputWrapper: {
    display: 'flex',
    margin: '1rem 0'
  },
  monitorInputWrapper: {
    display: 'flex'
  },
  inputText: {
    color: '#007bff'
  },
  endInput: {
    marginLeft: '0.8rem'
  },
  monitorWrapper: {
    justifyContent: 'space-between'
  }
}));

const GraphWrapper: FC<WrapperProps> = ({
  child,
  isMonitoring,
  routesChains,
  showDetails
}) => {
  const classes = useStyles();
  const d = moment().subtract(1, 'h');
  const [startTime, setStartTime] = useState(d);
  const [endTime, setEndTime] = useState(moment());
  const [startValue, setStartValue] = useState(d);
  const [endValue, setEndValue] = useState(moment());
  const [error, setError] = useState<string | null>(null);
  const handleStartDateClose = date => {
    const t = moment();
    if (startValue.isAfter(endTime) || startValue.isAfter(t)) {
      setStartValue(startTime);
      setError(
        'Start time cannot be greater than the end time.Kindly choose some other value'
      );
      setTimeout(() => {
        setError(null);
      }, 5000);
    } else {
      setError(null);
      setStartTime(startValue);
    }
  };

  const handleEndDateClose = date => {
    const t = moment();
    if (endValue.isBefore(startTime) || endValue.isAfter(t)) {
      setError(
        'End time cannot be less than the start time.Kindly choose some other value'
      );
      setTimeout(() => {
        setError(null);
      }, 5000);
      setEndValue(endTime);
    } else {
      setError(null);
      setEndTime(endValue);
    }
  };
  const handleStartChange = date => {
    setStartValue(date);
  };
  const handleEndChange = date => {
    setEndValue(date);
  };
  const valid = current => {
    return current.isBefore(moment());
  };
  const endTimestamp = endTime ? endTime.toISOString() : null;
  const dater = moment().subtract(1, 'h');
  const startTimestamp = startTime
    ? startTime.toISOString()
    : dater.toISOString();
  return (
    <div>
      {error && <Alert severity="error">{error}</Alert>}
      <div
        style={{ display: 'flex' }}
        className={isMonitoring ? classes.monitorWrapper : ''}
      >
        {isMonitoring && routesChains && showDetails && (
          <span
            style={{ display: 'flex', alignItems: 'center' }}
            onClick={() => showDetails(false, routesChains)}
          >
            <ArrowBackIcon color="primary" fontSize="large" />
            <span
              style={{
                fontSize: '1rem',
                fontWeight: 'bold',
                padding: '0 0.4rem',
                display: 'flex',
                alignItems: 'center'
              }}
            >
              {truncate(routesChains.name, 70)}
            </span>
          </span>
        )}
        <div
          className={
            isMonitoring ? classes.monitorInputWrapper : classes.inputWrapper
          }
        >
          <div>
            <div className={classes.inputText}>Start Time</div>
            <Datetime
              value={startValue}
              onClose={handleStartDateClose}
              isValidDate={valid}
              onChange={handleStartChange}
            />
          </div>
          <div className={classes.endInput}>
            <div className={classes.inputText}>End Time</div>
            <Datetime
              value={endValue}
              onClose={handleEndDateClose}
              isValidDate={valid}
              onChange={handleEndChange}
            />
          </div>
        </div>
      </div>
      {React.cloneElement(child, { startTimestamp, endTimestamp })}
    </div>
  );
};

export default GraphWrapper;
