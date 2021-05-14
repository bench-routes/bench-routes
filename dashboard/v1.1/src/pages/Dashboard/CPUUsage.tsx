import React, { FC, useContext } from 'react';
import Chart from 'react-apexcharts';
import { chartData } from '../../utils/queryTypes';
import { ThemeContext } from '../../layouts/BaseLayout';
import { ApexOptions } from 'apexcharts';

interface CPUUsageProps {
  cpuMetrics: chartData[];
}

const CPUUsage: FC<CPUUsageProps> = ({ cpuMetrics }) => {
  const themeMode = useContext(ThemeContext);
  let theme;
  if (themeMode === {}) {
    theme = 'light';
  } else {
    theme = themeMode;
  }
  const dataFormatted = cpuMetrics;
  const series = [
    {
      name: 'CPU',
      data: dataFormatted
    }
  ];
  const options: ApexOptions = {
    chart: {
      type: 'area'
    },
    xaxis: {
      type: 'category',
      categories: [],
      labels: {
        show: true,
        rotate: 45,
        rotateAlways: true,
        hideOverlappingLabels: true,
        trim: true
      }
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
      text: 'CPU usage (in percent)',
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
      mode: theme
    }
  };

  return (
    <>
      <Chart series={series} options={options} height="300" type="area" />
    </>
  );
};

export default React.memo(CPUUsage);
