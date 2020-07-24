import { rootRouteObject, headersObject, paramsObject, bodyObject, paramsTransformValue } from './types';

export const getRoutesMap = (data: rootRouteObject[]) => {
  let configRoutes = new Map();
    data.forEach((route: rootRouteObject) => {
      const uri = route.URL;
      if (!configRoutes.has(uri)) {
        configRoutes.set(uri, [
          {
            method: route.Method,
            body: route.Body,
            headers: route.Header,
            params: route.Params
          }
        ]);
      } else {
        configRoutes.set(uri, [
          ...configRoutes.get(uri),
          {
            method: route.Method,
            body: route.Body,
            headers: route.Header,
            params: route.Params
          }
        ]);
      }
    });
  return configRoutes;
}


export const populateParams = (params: paramsObject[] | bodyObject[] | headersObject[]) => {
  let arr: paramsTransformValue[] = [];
  if (params !== null && params !== undefined) {
    params.forEach(param => {
      arr.push({
        key: param.Name || param.OfType,
        value: param.Value
      });
    });
  }
  return arr;
}
