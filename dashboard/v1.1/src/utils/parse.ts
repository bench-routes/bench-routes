import stc from 'string-to-color';
import { rootRouteObject, headersObject, paramsObject, bodyObject, paramsTransformValue, LabelType } from './types';

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
            params: route.Params,
            labels: route.Labels
          }
        ]);
      } else {
        configRoutes.set(uri, [
          ...configRoutes.get(uri),
          {
            method: route.Method,
            body: route.Body,
            headers: route.Header,
            params: route.Params,
            labels: route.Labels
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

export const populateLabels = (labels: string[]) => {
  let labelArr: LabelType[] = [];
  labels.forEach((label: string) => {
    labelArr.push({
      name: label,
      color: stc(label)
    });
  });
  return labelArr;
}
