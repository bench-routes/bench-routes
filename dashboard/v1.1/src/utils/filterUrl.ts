// TODO: improved URL filtering
const filterUrl = (url: string) => {
    let res;
    if (url.startsWith('https://www.')) {
      res = url.replace('https://www.', '');
    } else if (url.startsWith('https://')) {
      res = url.replace('https://', '');
    } else if (url.startsWith('http://www.')) {
      res = url.replace('http://www.', '');
    } else if (url.startsWith('http://')) {
      res = url.replace('http://', '');
    }
    return res;
}

export {filterUrl}