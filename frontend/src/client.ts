import store from './store';

export async function baseClient(endpoint: string, opts?: RequestInit): Promise<Response> {
  return fetch(import.meta.env.VITE_BACKEND_URL + endpoint, {
    ...opts,
    headers: {
      ...opts?.headers,
      ...store.getters.headers
    }
  });
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