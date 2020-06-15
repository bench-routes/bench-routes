import React, { FC, useState } from 'react';
import Button from '@material-ui/core/Button';
import ButtonGroup from '@material-ui/core/ButtonGroup';
import DoneIcon from '@material-ui/icons/Done';
import { Badge } from 'reactstrap';

interface TypeProps {
  slice: string[];
  getRequestType(type: string);
}

const Type: FC<TypeProps> = ({ slice, getRequestType }) => {
  const [showDone, setShowDone] = useState(false);
  const [selectedRequest, setSelectedRequest] = useState('');
  const sendType = (type: string): void => {
    setShowDone(true);
    setSelectedRequest(type);
    getRequestType(type);
  };
  return (
    <div>
      <ButtonGroup
        color="primary"
        variant="contained"
        aria-label="outlined primary button group"
      >
        {slice.map((type: string) => (
          <Button onClick={() => sendType(type.toUpperCase())}>
            {type.toUpperCase()}
          </Button>
        ))}
      </ButtonGroup>
      {showDone ? (
        <>
          <DoneIcon
            style={{ marginLeft: '3%' }}
            fontSize="large"
            color="secondary"
          />{' '}
          <Badge color="success" style={{ fontSize: '13px', marginLeft: '1%' }}>
            {selectedRequest}
          </Badge>
        </>
      ) : null}
    </div>
  );
};

export default Type;
