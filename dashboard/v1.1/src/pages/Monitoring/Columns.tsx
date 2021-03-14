interface Column {
  label: string;
  minWidth: number;
  align: 'right' | 'center' | 'left';
  format?: (value: number) => string;
}

export const columns: Column[] = [
  { label: 'Name', minWidth: 170, align: 'left' },
  { label: 'Ping', minWidth: 100, align: 'center' },
  {
    label: 'Jitter',
    minWidth: 100,
    align: 'center',
    format: (value: number) => value.toLocaleString('en-US')
  },
  {
    label: 'Response time',
    minWidth: 100,
    align: 'center',
    format: (value: number) => value.toLocaleString('en-US')
  },
  {
    label: 'Response length',
    minWidth: 100,
    align: 'center',
    format: (value: number) => value.toFixed(2)
  },
  {
    label: 'Status',
    minWidth: 100,
    align: 'center',
    format: (value: number) => value.toFixed(2)
  },
  {
    label: '',
    minWidth: 10,
    align: 'center',
    format: (value: number) => value.toFixed(2)
  }
];
