import { ChartOptions, ChartValues } from '../layouts/Charts';

export const getChartOptions=(data)=>{
    console.log("inside getChartOptions")
    const yMin: ChartOptions[] = [];
            const yMean: ChartOptions[] = [];
            const yMax: ChartOptions[] = [];
            const yMdev: ChartOptions[] = [];
            const norTime: number[] = [];
            const timeStamp: string[] = [];
    
            if (data.length === 0) {
                // Probably send the required information
                // to the user via br-logger
                console.log('No data from the url');
              } else {
                let inst;
                for (inst of data) {
                  yMin.push(inst.Min);
                  yMean.push(inst.Mean);
                  yMax.push(inst.Max);
                  yMdev.push(inst.Mdev);
                  norTime.push(inst.relative);
                  timeStamp.push(inst.timestamp);
                }
              }
    
    const options: ChartOptions[] = [
        ChartValues(norTime, yMin, 'Minimum', 'rgba(75,192,192,0.4)'),
        ChartValues(norTime, yMean, 'Mean', 'rgba(75,192,2,0.4)'),
        ChartValues(norTime, yMax, 'Maximum', 'rgba(5,192,19,0.4)'),
        ChartValues(norTime, yMdev, 'Standard-Deviation', 'rgba(7,12,19,0.4)')
      ];

      return options;
}

 