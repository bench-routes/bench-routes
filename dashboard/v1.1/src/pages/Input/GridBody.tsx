import React, { FC, useState } from 'react';
import TextField from '@material-ui/core/TextField';

interface GridBodyProps {
  name: string;
}

interface pair {
  key: string;
  value: string;
}

const GridBody: FC<GridBodyProps> = ({ name }) => {
  const [header, setHeader] = useState<pair[]>([{ key: '', value: '' }]);
  const [items, setItems] = useState<number>(header.length);
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
    setItems(header.length);
  };
  return (
    <div>
      <h6>{name}</h6>
      <hr />
      <div className="row">
        <div className="col-md-8" style={{ display: 'inline-grid' }}>
          {/* Grid */}
          {header.map((_, index) => {
            return (
              <div className="row" style={{ margin: '3px' }}>
                <div className="col-md-6">
                  <TextField
                    id="outlined-basic"
                    size="small"
                    label="Key"
                    style={{ width: '100%' }}
                    variant="outlined"
                    onChange={e => updateItems(e.target.value, '', index)}
                  />
                </div>
                <div className="col-md-6">
                  <TextField
                    id="outlined-basic"
                    size="small"
                    style={{ width: '100%' }}
                    label="value"
                    variant="outlined"
                    onChange={e => updateItems('', e.target.value, index)}
                  />
                </div>
              </div>
            );
          })}
        </div>
        <div className="col-md-4">
          {/* Body */}
          <TextField
            id="outlined-multiline-flexible"
            label="JSONified"
            multiline
            rows={header.length * 3}
            value={bodyValue}
            variant="outlined"
            style={{ minHeight: '100%', width: '100%' }}
          />
        </div>
      </div>
    </div>
  );
};

export default GridBody;
