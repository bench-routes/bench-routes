import { useState, useEffect } from 'react';

export type APIResponse<T> = { status: string; data?: T };

export interface FetchState<T> {
  response: APIResponse<T>;
  error?: Error;
  isLoading: boolean;
}

export const useFetch = <T extends {}>(url: string, options?: RequestInit): FetchState<T> => {
  const [response, setResponse] = useState<APIResponse<T>>({ status: 'start fetching' });
  const [error, setError] = useState<Error>();
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const [requestSent, setRequestSent] = useState<boolean>(false);

  useEffect(() => {
    const fetchData = async () => {
      try {
        setIsLoading(true);
        setRequestSent(true);
        const res = await fetch(url, { cache: 'no-cache', credentials: 'same-origin', ...options });
        if (!res.ok) {
          setIsLoading(false);
          setRequestSent(false);
          throw new Error(res.statusText);
        }
        const json = (await res.json()) as APIResponse<T>;
        setResponse(json);
        setIsLoading(false);
      } catch (error) {
        setError(error);
      }
    };
    if (!requestSent) {
      setTimeout(() => {
        fetchData();
      }, 1000);
    }
  }, [url, options, requestSent]);
  return { response, error, isLoading };
};
