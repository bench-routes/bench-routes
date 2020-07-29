import React, { useEffect, useState } from 'react';
import Select from '@material-ui/core/Select';
import { FormControl, MenuItem, Typography } from '@material-ui/core';

interface DropdownItemType {
  value: string;
  name: string;
}

interface SelectProps {
  labelId?: string;
  id?: string;
  defaultValue?: string;
  label?: string;
  dropdownitems: DropdownItemType[];
  variant?: 'standard' | 'outlined' | 'filled';
  onSelection: (type: string) => void;
}

const MaterialSelect = (props: SelectProps) => {
  const [open, setOpen] = React.useState(false);
  const [dropdownList, setDropdownList] = useState<DropdownItemType[]>([]);
  const {
    variant,
    label,
    labelId,
    id,
    defaultValue,
    dropdownitems,
    onSelection
  } = props;
  useEffect(() => {
    if (dropdownitems.length === 0) {
      setDropdownList([
        {
          name: 'No data',
          value: 'No data'
        }
      ]);
    } else setDropdownList(dropdownitems);
  }, [dropdownitems]);
  const handleChange = event => {
    onSelection(event.target.value);
  };

  const handleClose = () => {
    setOpen(false);
  };

  const handleOpen = () => {
    setOpen(true);
  };

  return (
    <FormControl variant={variant} style={{ width: '100%' }}>
      <Typography style={{ paddingLeft: '0.5rem' }} variant="subtitle1">
        {label}
      </Typography>
      <Select
        labelId={labelId || 'label-id'}
        id={id || 'select-id'}
        open={open}
        onClose={handleClose}
        onOpen={handleOpen}
        onChange={handleChange}
        value={defaultValue}
      >
        {dropdownList.map((item: DropdownItemType) => (
          <MenuItem key={item.value} value={item.value}>
            {item.name}
          </MenuItem>
        ))}
      </Select>
    </FormControl>
  );
};

export default MaterialSelect;
