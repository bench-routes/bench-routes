import { useFetch } from './useFetch';
import { HOST_IP } from './types';
import { QueryResponse } from './queryTypes';

export interface APIResponse<T> {
  status: string;
  data: T;
}

export const init = (): QueryResponse => {
  return {
    evaluationTime: '',
    range: {
      start: 0,
      end: 0
    },
    timeSeriesPath: '',
    values: []
  };
};

const GetSystemData = () => {
  const { response, error } = useFetch(
    `${HOST_IP}/query?timeSeriesPath=storage/system`
  );

  if (response) return response;
  else return error;
};

const GetServiceState = () => {
  const { response, error } = useFetch(`${HOST_IP}/service-state`);

  if (response) return response;
  else return error;
};

export { GetSystemData, GetServiceState };
