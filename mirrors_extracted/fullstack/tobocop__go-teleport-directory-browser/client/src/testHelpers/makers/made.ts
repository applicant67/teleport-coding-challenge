export type Made<T> = { [P in keyof T]: jest.Mock | any };
