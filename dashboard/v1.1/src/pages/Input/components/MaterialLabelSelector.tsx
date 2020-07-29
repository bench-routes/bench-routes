import React, { Suspense, useEffect, useState } from 'react';
import { useTheme } from '@material-ui/core/styles';
import {
  HelpOutline as HelpOutlineIcon,
  Close as CloseIcon,
  Done as DoneIcon
} from '@material-ui/icons';
import Autocomplete, {
  AutocompleteCloseReason
} from '@material-ui/lab/Autocomplete';
import ButtonBase from '@material-ui/core/ButtonBase';
import InputBase from '@material-ui/core/InputBase';
import {
  CircularProgress,
  Popover,
  Tooltip,
  Typography
} from '@material-ui/core';
import { truncate } from '../../../utils/stringManipulations';
import { defaultLabel } from '../../../services/input';
import { HOST_IP, LabelType } from '../../../utils/types';
import stc from 'string-to-color';
import useStyles from '../styles/labels';

interface LabelSelectorType {
  updateLabels: (labels: LabelType[]) => void;
  defaultLabels: LabelType[];
}

export default function MaterialLabelSelector(props: LabelSelectorType) {
  const classes = useStyles();
  const { defaultLabels, updateLabels } = props;
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);
  const [value, setValue] = useState<LabelType[]>(defaultLabels);
  const [pendingValue, setPendingValue] = useState<LabelType[]>(defaultLabels);
  const [labelList, setLabelList] = useState<LabelType[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [labelValue, setLabelValue] = useState<string>('');
  const theme = useTheme();

  const handleAddLabel = () => {
    const newLabel = {
      name: labelValue,
      color: stc(labelValue)
    };
    setLabelList([...labelList, newLabel]);
  };

  useEffect(() => {
    fetch(`${HOST_IP}/get-labels`)
      .then(resp => resp.json())
      .then(response => {
        let labels: LabelType[] = [];
        const { data } = response;
        data.forEach(labelName => {
          labels.push({
            name: labelName,
            color: stc(labelName)
          });
        });
        return labels;
      })
      .then(labels => {
        setLabelList(labels);
        setValue(defaultLabels);
        setPendingValue(defaultLabels);
      })
      .then(() => setLoading(false));
  }, [loading, defaultLabels]);

  const handleClick = (event: React.MouseEvent<HTMLElement>) => {
    setPendingValue(value);
    setAnchorEl(event.currentTarget);
  };

  const handleClose = (
    event: React.ChangeEvent<{}>,
    reason: AutocompleteCloseReason
  ) => {
    if (reason === 'toggleInput') {
      return;
    }
    if (anchorEl) {
      anchorEl.focus();
    }
    setAnchorEl(null);
    setValue(pendingValue);
    updateLabels(pendingValue);
  };

  const open = Boolean(anchorEl);
  const id = open ? 'label' : undefined;
  const tooltipTitle = `Labels are a way to tag a group of routes together.
  Make sure not to add too many labels which may increase the 
  search space while loading the dashboard.`;

  return (
    <React.Fragment>
      <div className={classes.root}>
        <ButtonBase disableRipple aria-describedby={id} onClick={handleClick}>
          <Typography variant="subtitle1">Add Labels</Typography>
        </ButtonBase>
        <div className={classes.value}>
          {value.length > 0 ? (
            value.map((label: { name: string; color: string }) => (
              <div
                key={label.name}
                className={classes.tag}
                style={{
                  backgroundColor: label.color,
                  color: theme.palette.getContrastText(label.color)
                }}
              >
                {label.name}
              </div>
            ))
          ) : (
            <div
              key={defaultLabel.name}
              className={classes.tag}
              style={{
                backgroundColor: defaultLabel.color,
                color: theme.palette.getContrastText(defaultLabel.color)
              }}
            >
              {defaultLabel.name}
            </div>
          )}
        </div>
      </div>
      <Suspense fallback={<CircularProgress />}>
        <Popover id={id} open={open} anchorEl={anchorEl}>
          <div className={classes.header}>
            Labels
            <Tooltip title={tooltipTitle}>
              <HelpOutlineIcon color="action" />
            </Tooltip>
          </div>
          {loading ? (
            <CircularProgress />
          ) : (
            <Autocomplete
              open
              onClose={handleClose}
              multiple
              classes={{
                paper: classes.paper,
                option: classes.option,
                popperDisablePortal: classes.popperDisablePortal
              }}
              value={pendingValue}
              onChange={(event, newValue) => {
                setPendingValue(newValue);
              }}
              disableCloseOnSelect
              disablePortal
              renderTags={() => null}
              noOptionsText={
                <div className={classes.noOption} onMouseDown={handleAddLabel}>
                  Add new label
                  <div className={classes.tag}> {truncate(labelValue, 10)}</div>
                </div>
              }
              renderOption={(option, { selected }) => (
                <React.Fragment>
                  <DoneIcon
                    className={classes.iconSelected}
                    style={{ visibility: selected ? 'visible' : 'hidden' }}
                  />
                  <span
                    className={classes.color}
                    style={{ backgroundColor: option.color }}
                  />
                  <div className={classes.text}>
                    {truncate(option.name, 20)}
                  </div>
                  <CloseIcon
                    className={classes.close}
                    style={{ visibility: selected ? 'visible' : 'hidden' }}
                  />
                </React.Fragment>
              )}
              options={[...labelList].sort((a, b) => {
                // Display the selected labels first.
                let ai = value.map(el => el.name).indexOf(a.name);
                ai = ai === -1 ? value.length + labelList.indexOf(a) : ai;
                let bi = value.map(el => el.name).indexOf(b.name);
                bi = bi === -1 ? value.length + labelList.indexOf(b) : bi;
                return ai - bi;
              })}
              getOptionSelected={(option, value) => {
                if (option.name === value.name) {
                  return true;
                }
                return false;
              }}
              getOptionLabel={option => option.name}
              renderInput={params => (
                <InputBase
                  ref={params.InputProps.ref}
                  inputProps={params.inputProps}
                  autoFocus
                  className={classes.inputBase}
                  onChange={event => setLabelValue(event.target.value)}
                />
              )}
            />
          )}
        </Popover>
      </Suspense>
    </React.Fragment>
  );
}
