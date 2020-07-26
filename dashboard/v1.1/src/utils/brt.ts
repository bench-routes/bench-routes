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

export const round = (n: string): number => {
	const num = parseInt(n, 10);
	return Math.round(num * 10) / 10;
};

export const formatResLength = (resLength: number): string => {
	if (resLength > 9999999999) {
		return (resLength / 1000000000).toFixed(2).toString() + 'B';
	} else if (resLength > 999999) {
		return (resLength / 1000000).toFixed(2).toString() + 'M';
	} else if (resLength > 999) {
		return (resLength / 1000).toFixed(2).toString() + 'K';
	}
	return resLength.toString();
};

export default TimeInstance;