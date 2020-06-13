import React, { FC, useState } from 'react';
import Button from '@material-ui/core/Button';
import ButtonGroup from '@material-ui/core/ButtonGroup';
import { makeStyles, createStyles, Theme } from '@material-ui/core/styles';
import DoneIcon from '@material-ui/icons/Done';
import { Badge } from 'reactstrap';

const useStyles = makeStyles((theme: Theme) =>
  createStyles({
    root: {
      display: 'flex',
      flexDirection: 'column',
      '& > *': {
        margin: theme.spacing(1)
      }
    }
  })
);

interface BlockProps {
  name: string;
}

const Block: FC<BlockProps> = ({ name }) => (
  <Button>{name.toUpperCase()}</Button>
);

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
