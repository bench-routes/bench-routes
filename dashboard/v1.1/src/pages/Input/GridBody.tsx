import React, { FC, useState } from 'react';
import TextField from '@material-ui/core/TextField';

interface GridBodyProps {
  name: string;
  updateParent: React.Dispatch<React.SetStateAction<pair[] | undefined>>;
}

export interface pair {
  key: string;
  value: string;
}

const GridBody: FC<GridBodyProps> = ({ name, updateParent }) => {
  const [header, setHeader] = useState<pair[]>([{ key: '', value: '' }]);
  const [bodyValue, setBodyValue] = useState<string>('');
  const updateItems = (key: string, value: string, i: number) => {
    if (key === '') {
      header[i].value = value;
    }
    if (value === '') {
      header[i].key = key;
    }
    if (i === header.length - 1) {
      header.push({ key: '', value: '' });
    }
    setBodyValue(JSON.stringify(header, null, 4));
    updateParent(header);
  };
  const updateBody = (content: string) => {
    setBodyValue(content);
    const inJSON = JSON.parse(content) as pair[];
    header.length = 0;
    for (const pair of inJSON) {
      header.push({
        key: pair.key,
        value: pair.value
      });
    }
    setHeader(header);
    updateParent(header);
  };
  return (
    <div>
      <h6 style={{ fontWeight: 'bold' }}>{name}</h6>
      <hr />
      <div className="row">
        <div className="col-md-8" style={{ display: 'inline-grid' }}>
          {header.map((head, index) => (
            <div className="row" style={{ margin: '3px' }} key={index}>
              <div className="col-md-6">
                <TextField
                  id="outlined-basic"
                  size="small"
                  label="Key"
                  value={head.key}
                  style={{ width: '100%' }}
                  variant="outlined"
                  onChange={e => updateItems(e.target.value, '', index)}
                />
              </div>
              <div className="col-md-6">
                <TextField
                  id="outlined-basic"
                  size="small"
                  value={head.value}
                  style={{ width: '100%' }}
                  label="value"
                  variant="outlined"
                  onChange={e => updateItems('', e.target.value, index)}
                />
              </div>
            </div>
          ))}
        </div>
        <div className="col-md-4">
          <TextField
            id="outlined-multiline-flexible"
            label="JSONified"
            multiline
            rows={header.length * 3}
            value={bodyValue}
            variant="outlined"
            onChange={e => updateBody(e.target.value)}
            style={{
              minHeight: '100%',
              width: '100%',
              backgroundColor: '#f9f9f9'
            }}
          />
        </div>
      </div>
    </div>
  );
};

export default GridBody;
