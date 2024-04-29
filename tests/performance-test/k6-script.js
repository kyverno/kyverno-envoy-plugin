import http from 'k6/http';
import { check, group, sleep } from 'k6';

export const options = {
  stages: [
    { duration: '30s', target: 100 }, // Ramp-up to 100 virtual users over 30 seconds
    { duration: '1m', target: 100 }, // Stay at 100 virtual users for 1 minute
    { duration: '30s', target: 0 }, // Ramp-down to 0 virtual users over 30 seconds
  ],
};

const BASE_URL = 'minikube ip with sample application'; // Replace with your application URL

export default function () {
  group('GET /book with admin token', () => {
    const res = http.get(`${BASE_URL}/book`, {
      headers: { 'Authorization': 'Bearer your_admin_token' },
    });
    check(res, {
      'is status 200': (r) => r.status === 200,
    });
  });

  group('GET /book with guest token', () => {
    const res = http.get(`${BASE_URL}/book`, {
      headers: { 'Authorization': 'Bearer your_guest_token' },
    });
    check(res, {
      'is status 200': (r) => r.status === 200,
    });
  });

  group('POST /book with guest token', () => {
    const res = http.post(`${BASE_URL}/book`, {
      headers: { 'Authorization': 'Bearer your_guest_token' },
    });
    check(res, {
      'is status 403': (r) => r.status === 403,
    });
  });

  sleep(1); // Sleep for 1 second between iterations
}