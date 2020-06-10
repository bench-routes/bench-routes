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
    datalabels: {
      enabled: false
    },
    stroke: {
      show: true,
      curve: 'straight',
      lineCap: 'butt',
      width: 3
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
    datalabels: {
      enabled: false
    },
    stroke: {
      show: true,
      curve: 'straight',
      lineCap: 'butt',
      width: 2
    },
    subtitle: {
      text: 'Cache (in bytes)',
      align: 'center'
    }
  };

  return (
    <div className="row">
      <div className="col-md-6">
        <Chart series={seriesDiskIO} options={optionsDiskIO} height="300" />
      </div>
      <div className="col-md-6">
        <Chart series={seriesCache} options={optionsCache} height="300" />
      </div>
    </div>
  );
};

export default React.memo(DiskUsage);
