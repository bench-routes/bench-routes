import { useFetch } from './useFetch';
import { HOST_IP } from './types';

export interface APIResponse<T> {
  status: string;
  data: T;
}

export const init = (): {status: string, data: any} => {
  return {
    status: '',
    data: {}
  }
}

const GetSystemData = () => {

  const { response, error } = useFetch(
    `${HOST_IP}/query?timeSeriesPath=storage/system`
  );

  if (response) return response
  else return error

};

const GetServiceState = () => {

  const { response, error } = useFetch(
    `${HOST_IP}/service-state`
  );

  if (response) return response
  else return error

};



export { GetSystemData, GetServiceState }