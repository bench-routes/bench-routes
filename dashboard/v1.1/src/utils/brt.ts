enum TimeInstance {
	Microsecond = 1000,
	Millisecond = 1000000,
	Second = 1000000000,
	Minute = 60000000000,
	Hour = 60 * 60000000000,
	Day = 24 * 60 * 60000000000,
};

export const formatTime = (data: string) => {
  const times = data.split('|');
  const l = times.length;
  return `${times[l-3]}:${times[l-2]}:${times[l-1]}`;
};

export default TimeInstance;