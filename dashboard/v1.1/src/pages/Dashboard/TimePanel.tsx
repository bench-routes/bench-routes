import { Badge } from 'reactstrap';
import React from 'react';

const TimePanel = ({ fetchTime, evaluationTime }) => {
  return (
    <div style={{ borderLeft: '3px solid #2196f3' }}>
      {fetchTime ? (
        <h6
          style={{
            margin: 0,
            padding: 10,
            display: 'flex',
            alignItems: 'center'
          }}
        >
          Query fetched in
          <Badge
            color="success"
            style={{ fontSize: '0.97rem', marginLeft: 10 }}
          >
            {fetchTime.toFixed(3)}ms
          </Badge>
        </h6>
      ) : null}

      {evaluationTime !== '' ? (
        <h6
          style={{
            margin: 0,
            padding: 10,
            display: 'flex',
            alignItems: 'center'
          }}
        >
          Query Time :
          <Badge
            color="success"
            style={{ fontSize: '0.97rem', marginLeft: 10 }}
          >
            {evaluationTime}
          </Badge>
        </h6>
      ) : null}
    </div>
  );
};

export default TimePanel;
