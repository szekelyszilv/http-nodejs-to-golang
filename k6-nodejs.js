import http from 'k6/http';
import { check } from 'k6';

export const options = {
    vus: 3,
    duration: '10s',
};

export default () => {
    const urlRes = http.get('http://localhost:3000/nodejs');
    check(urlRes, { 'status returned 200': (r) => r.status == 200 })
};
