import BRConnect from "../../utils/connection";

export const opts = (operation: string, connection : BRConnect, type: string): void => {
    switch (operation) {
      case 'start':
          if(type==='ping')
          {
            
              connection.signalPingStart().then(res => {
                if (res.data) {
                  alert('Ping routine started');
                }
              });
          }
          else if(type==='jitter')
          {
              
            connection.signalJitterStart().then(res => {
                if (res.data) {
                  alert('Jitter routine started');
                }
              });
          }
          else if(type==='floodPing')
          {
             
            connection.signalFloodPingStart().then(res => {
                if (res.data) {
                  alert('Flood Ping routine started');
                }
              });
          }
          else{
              alert("Something went wrong");
          }
        break;
      case 'stop':
          if(type==='ping')
          {

              connection.signalPingStop().then(res => {
                if (res.data) {
                  alert('Ping routine stopped');
                }
              });
          }
          else if(type==='jitter')
          {
            connection.signalJitterStop().then(res => {
                if (res.data) {
                  alert('Jitter routine stopped');
                }
              });
          }
          else if (type==='floodPing')
          {
            connection.signalFloodPingStop().then(res => {
                if (res.data) {
                  alert('Flood Ping routine stopped');
                }
              });
          }
          else{
              alert("Something went wrong");
          }
        break;
    }
  };