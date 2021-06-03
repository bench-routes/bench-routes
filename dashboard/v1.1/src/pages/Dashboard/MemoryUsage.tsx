import React, { FC, useContext } from 'react';
import Chart from 'react-apexcharts';
import { chartData } from '../../utils/queryTypes';
import { ThemeContext } from '../../layouts/BaseLayout';

interface MemoryUsagePercentProps {
  memoryUsagePercentMetrics: chartData[];
}

const MemoryUsagePercent: FC<MemoryUsagePercentProps> = ({
  memoryUsagePercentMetrics
}) => {
  const themeMode = useContext(ThemeContext);
  const dataFormatted = memoryUsagePercentMetrics;
  const series = [
    {
      name: 'Memory',
      data: dataFormatted
    }
  ];
  const options = {
    chart: {
      type: 'area'
    },
    dataLabels: {
      enabled: false
    },
    stroke: {
      show: true,
      curve: 'straight',
      lineCap: 'butt',
      width: 1
    },
    subtitle: {
      text: 'Memory (RAM) usage  (in percent)',
      align: 'center'
    },
    fill: {
      opacity: 1,
      type: 'gradient',
      gradient: {
        shade: 'dark',
        type: 'vertical',
        shadeIntensity: 0.6,
        inverseColors: true,
        opacityFrom: 1,
        opacityTo: 0.2
      }
    },
    theme: {
      mode: themeMode
    }
  };

  return (
    <>
      <Chart series={series} options={options} height="300" type="area" />
    </>
  );
};

export default React.memo(MemoryUsagePercent);
