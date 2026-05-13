export function setLocalStore(store) {
  localStorage.setItem('store', JSON.stringify(store));
}

export function getLocalStore() {
  return JSON.parse(localStorage.getItem('store'));
}
