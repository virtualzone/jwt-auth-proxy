interface AjaxResult {
  json: any
  status: number
  objectId: string
}

export default class Ajax {
  static URL_PREFIX: string = "http://localhost:8080";

  static async query(method: string, url: string, data?: any): Promise<AjaxResult> {
    if (process.env.NODE_ENV === 'development') {
      url = Ajax.URL_PREFIX + url;
    }
    let headers = new Headers();
    if (window.sessionStorage.getItem("accessToken") != null) {
      headers.append("Authorization", "Bearer " + window.sessionStorage.getItem("accessToken"))
    }
    if (data) {
      headers.append("Content-Type", "application/json");
    }
    let options: RequestInit = {
      method: method,
      mode: "cors",
      cache: "no-cache",
      credentials: "same-origin",
      headers: headers
    };
    if (data) {
      options.body = JSON.stringify(data);
    }
    return new Promise<AjaxResult>(function(resolve, reject) {
      fetch(url, options).then((response) => {
        response.json().then(json => {
          resolve({
            json: json,
            status: response.status,
            objectId: response.headers.get("X-Object-Id")
          } as AjaxResult);
        }).catch(err => {
          resolve({
            json: {},
            status: response.status,
            objectId: response.headers.get("X-Object-Id")
          } as AjaxResult);
        });
      }).catch(err => {
        reject(err);
      });
    });
  }

  static async postData(url: string, data?: any): Promise<AjaxResult> {
    return Ajax.query("POST", url, data);
  }

  static async get(url: string): Promise<AjaxResult> {
    return Ajax.query("GET", url);
  }
}
