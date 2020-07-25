import { HOST_IP, routeOptionsInterface } from '../../utils/types';
import { getRoutesMap } from '../../utils/parse';

export interface TableRouteType {
  route: string;
  methods: string[];
}

export interface IntervalType {
  test: string;
  duration: number;
  unit: string;
}

export const onRowDelete = (oldData: TableRouteType, setConfigRoutes) => {
  fetch(`${HOST_IP}/delete-route`, {
    method: 'post',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      actualRoute: oldData.route
    })
  })
    .then(resp => resp.json())
    .then(response => {
      const { data } = response;
      let configRoutes: Map<
        string,
        routeOptionsInterface[]
      > = getRoutesMap(data);
      setConfigRoutes(configRoutes);
    }, err => {
      console.error(err);
    });
};
