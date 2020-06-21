import React from 'react';
import { tableIcons } from '../../../utils/tableIcons';

import MaterialTable from 'material-table';

const SearchTable = (props: any) => {
  return (
    <MaterialTable
      icons={tableIcons}
      style={{ marginTop: '10vh' }}
      options={{
        headerStyle: {
          fontSize: 18,
          fontWeight: 'normal'
        }
      }}
      {...props}
    />
  );
};

export default SearchTable;
