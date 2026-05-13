import { getLocalStore } from '../localStorage';

function checkStatus(response) {
  if (response.status >= 200 && response.status < 300) {
    return response;
  }
  const error = new Error(response.statusText);
  error.response = response;
  throw error;
}

function parseJSON(response) {
  return response.json();
}

function getBearerTokenHeader() {
  const store = getLocalStore();
  const headers = {};
  if (store) {
    headers.Authorization = `Bearer ${store.sessionID}`;
  }
  return headers;
}

export default {
  async post(route, body) {
    const response = await fetch(`api${route}`, {
      method: 'POST',
      headers: {
        Accept: 'application/json',
        'Content-Type': 'application/json',
        ...getBearerTokenHeader(),
      },
      body: JSON.stringify(body),
    });
    const responseChecked = await checkStatus(response);
    return parseJSON(responseChecked);
  },
  async delete(route) {
    const response = await fetch(`api${route}`, {
      method: 'DELETE',
      headers: getBearerTokenHeader(),
    });
    return checkStatus(response);
  },
  async get(route) {
    const response = await fetch(`api${route}`, {
      method: 'GET',
      headers: {
        Accept: 'application/json',
        ...getBearerTokenHeader(),
      },
    });
    const responseChecked = await checkStatus(response);
    return parseJSON(responseChecked);
  },
  async patch(route) {
    const response = await fetch(`api${route}`, {
      method: 'PATCH',
      headers: {
        Accept: 'application/json',
        ...getBearerTokenHeader(),
      },
    });
    const responseChecked = await checkStatus(response);
    return parseJSON(responseChecked);
  },
};
