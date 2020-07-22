import React from 'react';
import { tableIcons } from '../../../utils/tableIcons';

import MaterialTable from 'material-table';
import { Chip } from '@material-ui/core';

export interface TableRowData {
  methods: string[];
  route: string;
  tableData: { id: number };
}

export const columns = [
  {
    field: 'route',
    title: 'Route'
  },
  {
    field: 'methods',
    title: 'Methods',
    render: (rowData: TableRowData) =>
      rowData.methods.map(m => (
        <Chip
          key={rowData.route + m}
          variant="outlined"
          color="primary"
          label={m}
          style={{ marginLeft: '2px' }}
        />
      ))
  }
];

const SearchTable = (props: any) => {
  return (
    <MaterialTable
      icons={tableIcons}
      style={{ marginTop: '2vh' }}
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
