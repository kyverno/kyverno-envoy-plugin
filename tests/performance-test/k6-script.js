import http from 'k6/http';
import { check, group, sleep } from 'k6';

export const options = {
  stages: [
    { duration: '30s', target: 100 }, // Ramp-up to 100 virtual users over 30 seconds
    { duration: '1m', target: 100 }, // Stay at 100 virtual users for 1 minute
    { duration: '30s', target: 0 }, // Ramp-down to 0 virtual users over 30 seconds
  ],
};

/*
Replace ip for every scenerio generate the URL with these commands 

echo SERVICE_PORT=$(kubectl -n demo get service testapp -o jsonpath='{.spec.ports[?(@.port==8080)].nodePort}')
echo SERVICE_HOST=$(minikube ip)
echo SERVICE_URL=$SERVICE_HOST:$SERVICE_PORT
echo $SERVICE_URL

http://192.168.49.2:31541

*/

const BASE_URL = 'http://192.168.49.2:31700'; // Replace with your application URL 

export default function () {
  group('GET /book with guest token', () => {
    const res = http.get(`${BASE_URL}/book`, {
      headers: { 'Authorization': 'Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjIyNDEwODE1MzksIm5iZiI6MTUxNDg1MTEzOSwicm9sZSI6Imd1ZXN0Iiwic3ViIjoiWVd4cFkyVT0ifQ.ja1bgvIt47393ba_WbSBm35NrUhdxM4mOVQN8iXz8lk' },
    });
    check(res, {
      'is status 200': (r) => r.status === 200,
    });
  });

  sleep(1); // Sleep for 1 second between iterations
}
