import React from 'react';
import { tableIcons } from '../../../utils/tableIcons';

import MaterialTable from 'material-table';
import { Chip } from '@material-ui/core';
import { truncate } from '../../../utils/stringManipulations';

export interface TableRowData {
  methods: string[];
  route: string;
  tableData: { id: number };
}

export const columns = [
  {
    field: 'route',
    title: 'Route',
    render: (rowData: any) => <>{truncate(rowData.route, 40)}</>
  },
  {
    field: 'methods',
    title: 'Methods',
    render: (rowData: TableRowData) =>
      rowData.methods.map((m, i) =>
        i % 2 ? (
          <Chip
            key={rowData.route + m}
            color="secondary"
            label={m}
            size="small"
            style={{ marginLeft: '1%' }}
          />
        ) : (
          <Chip
            key={rowData.route + m}
            color="primary"
            label={m}
            size="small"
            style={{ marginLeft: '1%' }}
          />
        )
      )
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
