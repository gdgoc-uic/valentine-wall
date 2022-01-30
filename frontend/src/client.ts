import store from './store';

const backendUrl = import.meta.env.VITE_BACKEND_URL;

interface BaseResponse {
  rawResponse: Response
}

export interface APIResponse extends BaseResponse {
  data: any
}

export class APIResponseError extends Error implements BaseResponse {
  rawResponse: Response;

  constructor(message: string, rawResponse: Response) {
    super(message);
    this.rawResponse = rawResponse;
  }

  static fromResponseWithJson(resp: Response, data?: any): APIResponseError {
    let errorMessage = 'Something went wrong.';
    if (data) {
      if ('error_message' in data) {
        errorMessage = data['error_message'];
      } else if ('messsage' in data) {
        errorMessage = data['message'];
      }
    }
    return new APIResponseError(errorMessage, resp);
  }
}

export function expandAPIEndpoint(endpoint: string): string {
  let finalEndpoint = endpoint;
  if (backendUrl[backendUrl.length - 1] == '/' && finalEndpoint[0] == '/') {
    finalEndpoint = endpoint.substring(1);
  }
  return backendUrl + finalEndpoint;
}

export async function baseClient(endpoint: string, opts?: RequestInit): Promise<APIResponse> {
  const resp = await fetch(expandAPIEndpoint(endpoint), {
    ...opts,
    headers: {
      ...opts?.headers,
      ...store.getters.headers
    }
  });

  const contentType = resp.headers.get('Content-Type');
  if (contentType && contentType == 'application/json') {
    const json = await resp.json();
    if (!resp.ok) {
      throw APIResponseError.fromResponseWithJson(resp, json);
    }

    return {
      rawResponse: resp,
      data: json
    }
  }

  const data = await resp.text();
  return {
    rawResponse: resp,
    data
  };
}

const client = {
  get(endpoint: string, opts?: RequestInit) {
    return baseClient(endpoint, { method: 'GET', ...opts });
  },
  post(endpoint: string, opts?: RequestInit) {
    return baseClient(endpoint, { method: 'POST', ...opts });
  },
  postJson(endpoint: string, payload: any, opts?: RequestInit) {
    return this.post(endpoint, {
      ...opts,
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(payload)
    });
  },
  delete(endpoint: string, opts?: RequestInit) {
    return baseClient(endpoint, { method: 'DELETE', ...opts });
  }
}

export default client;