import {
  HOST_IP,
  rootRouteObject
} from '../utils/types';

interface interval {
  Test: string;
  Duration: string;
  Type: string;
}

export const fetchConfigIntervals = async (setConfigIntervals) => {
  try {
    const response = await fetch(`${HOST_IP}/get-config-intervals`).then(
      resp => {
        return resp.json();
      }
    );
    const { data } = response;
    let intervals: any = [];
    data.forEach((interval: interval) => {
      intervals.push({
        test: interval.Test,
        duration: interval.Duration,
        unit: interval.Type
      });
    });
    setConfigIntervals(intervals);
  } catch (e) {
    console.log(e);
  }
};

export const fetchConfigRoutes = async (setConfigRoutes) => {
  const response = await fetch(`${HOST_IP}/get-config-routes`).then(resp => {
    return resp.json();
  });
  const { data } = response;
  let configRoutes = new Map();
  data.forEach((route: rootRouteObject) => {
    const uri = route.URL
    if (!configRoutes.has(uri)) {
      configRoutes.set(uri, [{
          method: route.Method,
          body: route.Body,
          headers: route.Header,
          params: route.Params
        }]);
    } else {
      configRoutes.set(uri, [
        ...configRoutes.get(uri), {
          method: route.Method,
          body: route.Body,
          headers: route.Header,
          params: route.Params
        }]);
    }
  });
  setConfigRoutes(configRoutes);
};
