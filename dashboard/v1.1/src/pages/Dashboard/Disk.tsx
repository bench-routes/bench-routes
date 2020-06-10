import React, { FC } from 'react';
import Chart from 'react-apexcharts';
import { chartData } from '../../utils/queryTypes';

interface DiskUsageProps {
  diskIO: chartData[];
  cache: chartData[];
}

const DiskUsage: FC<DiskUsageProps> = ({ diskIO, cache }) => {
  const seriesDiskIO = [
    {
      name: 'Disk IO in bytes (+ve means write / -ve means read)',
      data: diskIO
    }
  ];
  const seriesCache = [
    {
      name: 'Cache (in bytes)',
      data: cache
    }
  ];

  const optionsDiskIO = {
    chart: {
      type: 'area',
      background: '#fff'
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
      text: 'Disk IO in bytes (+ve means write / -ve means read)',
      align: 'center'
    }
  };
  const optionsCache = {
    chart: {
      type: 'area',
      animations: {
        enabled: true,
        easing: 'easeinout',
        speed: 800,
        animateGradually: {
          enabled: true,
          delay: 150
        },
        dynamicAnimation: {
          enabled: true,
          speed: 350
        }
      },
      background: '#fff'
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
      text: 'Cache (in bytes)',
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
    }
  };

  return (
    <div className="row">
      <div className="col-md-6">
        <Chart
          series={seriesDiskIO}
          options={optionsDiskIO}
          height="300"
          type="area"
        />
      </div>
      <div className="col-md-6">
        <Chart
          series={seriesCache}
          options={optionsCache}
          height="300"
          type="area"
        />
      </div>
    </div>
  );
};

export default React.memo(DiskUsage);
