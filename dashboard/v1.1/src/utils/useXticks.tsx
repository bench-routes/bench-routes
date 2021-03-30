import React, { useState, useEffect, useContext, createContext } from 'react';

const XTICKS_TOKEN = 'xticks';
const DEFAULT_XTICKS = '10';

function useProvideXticks() {
  const [xticks, setXticks] = useState(DEFAULT_XTICKS);

  const updateXticks = (value: string) => {
    setXticks(value);
    localStorage.setItem(XTICKS_TOKEN, value);
  };

  useEffect(() => {
    const item = localStorage.getItem(XTICKS_TOKEN);
    updateXticks(item ? item : DEFAULT_XTICKS);
  }, []);

  // Return the xticks and updateXticks method
  return {
    xticks,
    updateXticks
  };
}

const XticksContext = createContext({});
export function ProvideXticks({
  children
}: {
  children?: JSX.Element;
}): JSX.Element {
  const xticks = useProvideXticks();
  return (
    <XticksContext.Provider value={xticks}>{children}</XticksContext.Provider>
  );
}

export const useXticks = () => useContext(XticksContext);
export default useXticks;
