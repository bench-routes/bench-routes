export interface RouteFetchAll {
  url: string;
}

interface StoreType {
    routeDetails: any;
}

export default class BRConnect {
  public store: StoreType;
  private socketConn: WebSocket;
  private urlSocketConn: string;

  constructor() {
    this.store = {
      routeDetails: {}
    };
    this.urlSocketConn = 'ws://localhost:9090/websocket';
    this.socketConn = new WebSocket(this.urlSocketConn);
    this.socketConn.onopen = () => {
      this.socketConn.send('hi from br-e');

      // initialise connection
      this.routeDetails();
    };
  }

  public routeDetails(): Promise<any> {
    return this.sendMessage('route-details');
  }

  public signalPingStart(): Promise<any> {
    return this.sendMessage('force-start-ping');
  }

  public signalPingStop(): Promise<any> {
    return this.sendMessage('force-stop-ping');
  }

  public signalPingRouteFetchAllTimeSeries(route: string): Promise<any> {
    const inst: RouteFetchAll = {
      url: route
    };
    return this.sendAndReceiveMessage('Qping-route ' + JSON.stringify(inst));
  }

  public signalJitterRouteFetchAllTimeSeries(route: string): Promise<any> {
    const inst: RouteFetchAll = {
      url: route
    };
    return this.sendAndReceiveMessage('Qjitter-route ' + JSON.stringify(inst));
  }

  public signalFloodPingRouteFetchAllTimeSeries(route: string): Promise<any> {
    const inst: RouteFetchAll = {
      url: route
    };
    return this.sendAndReceiveMessage('Qflood-ping-route ' + JSON.stringify(inst));
  }

  public signalReqResDelayRouteFetchAllTimeSeries(route: string): Promise<any> {
    const inst: RouteFetchAll = {
      url: route
    };
    return this.sendAndReceiveMessage('Qrequest-response-delay ' + JSON.stringify(inst));
  }

  private sendAndReceiveMessage(message: string): Promise<any> {
    return new Promise((res: any, rej: any) => {
      this.socketConn.send(message);
      this.socketConn.onmessage = (m: any) => {
        res(m);
        console.log("M:::", m);
      };
      this.socketConn.onerror = (e: any) => {
        rej(e);
      };
  });
  }

  private sendMessage(message: string): Promise<any> {
    return new Promise((res: any, rej: any) => {
      if (this.socketConn.CONNECTING !== 0) {
        this.socketConn.send(message);
        this.socketConn.onmessage = (m: any) => {
          const data: string = m.data;
          const dataJSON: object = JSON.parse(data);
          res(dataJSON);
        };
        this.socketConn.onerror = (e: any) => {
          rej(e);
        };
      } else {
        this.socketConn.onopen = () => {
          this.socketConn.send('hi from br-e2');
          this.socketConn.send(message);
          this.socketConn.onmessage = (m: any) => {
              const data: string = m.data;
              const dataJSON: object = JSON.parse(data);
              res(dataJSON);
          };
          this.socketConn.onerror = (e: any) => {
              rej(e);
          };
        };
      }
    });
  }
}
