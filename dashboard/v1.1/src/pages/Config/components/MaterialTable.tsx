import React, { useState } from 'react';
import { tableIcons } from '../../../utils/tableIcons';

import MaterialTable from 'material-table';

const SearchTable = (props: any) => {
  const [selectedRow, setSelectedRow] = useState(null);

  return (
    <MaterialTable
      icons={tableIcons}
      style={{ marginTop: '10vh' }}
      onRowClick={(evt, selectedRow: any) =>
        setSelectedRow(selectedRow?.tableData.id)
      }
      options={{
        headerStyle: {
          fontSize: 20
        },
        rowStyle: rowData => ({
          backgroundColor:
            selectedRow === rowData.tableData.id ? '#EEE' : '#FFF'
        })
      }}
      {...props}
    />
  );
};

export default SearchTable;
