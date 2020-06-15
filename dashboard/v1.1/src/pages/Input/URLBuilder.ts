import { pair } from './GridBody';

export default class URLUtils {
  static formatParams(params: pair[] | undefined): string {
    if (params === undefined) {
      return '';
    }
    let p: string = '';
    for (const param of params) {
      p += `${param.key}=${param.value}&`;
    }
    return p.substr(0, p.length - 1);
  }
} 

export class URLBuilder {
  private url: string;
  private headers: pair[] | undefined;
  private params: pair[] | undefined;
  private body: pair[] | undefined;
  constructor(url: string, headers: pair[] | undefined, params: pair[] | undefined, body: pair[] | undefined) {
    this.url = url;
    this.headers = headers;
    this.params = params;
    this.body = body;
  }

  send(type: string): {success: boolean, status: number, response: string} {
    switch(type) {
      case 'get':
        console.warn('sending')
        const params = this.formatParams(this.params);
        const headers = this.formatHeaders(this.headers);
        const body = this.formatBody(this.body);
        let url: string = '';
        if (params.length === 0) {
          url = this.url;
        } else {
          url = `${this.url}?${params}`;
        }
        console.warn('url is ', url)
        fetch(url, {
          method: 'get',
          // mode:'no-cors',
          headers: headers
        }).then(res => res.json())
          .then((response) => {
            console.warn('got response as ', response)
          });
    }
    return {
      success: false,
      status: 400,
      response: '',
    };
  }

  formatParams(params: pair[] | undefined): string {
    if (params === undefined) {
      return '';
    }
    let p: string = '';
    for (const param of params) {
      p += `${param.key}=${param.value}&`;
    }
    return p.substr(0, p.length - 1);
  }

  formatHeaders(headers: pair[] | undefined): Record<string, string> {
    if (headers === undefined) {
      return {};
    }
    let instance: Record<string, string> = {};
    for (const header of headers) {
      instance[header.key] = header.value;
    }
    return instance;
  }

  formatBody(body: pair[] | undefined): Record<string, string> {
    if (body === undefined) {
      return {};
    }
    let instance: Record<string, string> = {};
    for (const b of body) {
      instance[b.key] = b.value;
    }
    return instance;
  }
}