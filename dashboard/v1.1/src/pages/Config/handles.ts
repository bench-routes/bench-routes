import { HOST_IP, routeOptionsInterface } from '../../utils/types';
import { getRoutesMap } from '../../utils/parse';

export interface TableRouteType {
  route: string;
  methods: string[];
  tableData?: {
    id: number;
    editing?: string | undefined;
  };
}

export interface IntervalType {
  test: string;
  duration: number;
  unit: string;
}

export const handleRowDelete = async (
  oldData: TableRouteType,
  setConfigRoutes: React.Dispatch<
    React.SetStateAction<Map<string, routeOptionsInterface[]>>
  >,
  tableData: TableRouteType[],
  setTableData: React.Dispatch<React.SetStateAction<TableRouteType[]>>
) => {
  let response;

  try {
    await fetch(`${HOST_IP}/delete-route`, {
      method: 'post',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        actualRoute: oldData.route
      })
    });
  } catch (err) {
    console.log(err);
  }

  const dataDelete = [...tableData];
  const index = oldData.tableData?.id;
  dataDelete.splice(index!, 1);
  setTableData([...dataDelete]);

  const { data } = await response.json();
  const configRoutes: Map<string, routeOptionsInterface[]> = getRoutesMap(data);
  setConfigRoutes(configRoutes);
};
