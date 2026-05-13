export interface ApiClient {
  authenticate(username: string, password: string): Promise<boolean>
}

// TODO: Add CSRF protection
export class ApiClientImpl {
  private baseRoute = '/api';

  authenticate(username: string, password: string): Promise<boolean> {
    return fetch(`${this.baseRoute}/authenticate`, {
      method: 'POST',
      cache: 'no-cache',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        username,
        password,
      }),
    }).then((r) => {
      if (r.status === 204) {
        return true;
      }
      // TODO: Should include better error handling and return status code
      throw new Error(`${r.status} Failed code`);
    });
  }
}
