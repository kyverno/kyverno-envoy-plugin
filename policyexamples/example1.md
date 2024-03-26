# Example 

Policy example where the user name is encoded in base64 authentication (Basic http token)

- Here user alice is granted a guest role and can perform a GET request to /productpage.
- And user bob is granted an admin role and can perform a GET to /productpage and /api/v1/products.

```
base64_decode(YWxpY2U6cGFzc3dvcmQ=)   =  alice:password
base64_decode(Ym9iOnBhc3N3b3Jk)       =  bob:password
```


Below is the example of good request payload which should pass the policy 
```json
{
    "attributes": {
        "request": {
            "http": {
                "method": "GET",
                "path": "/productpage",
                "headers": {
                    "authorization": "Basic YWxpY2U6cGFzc3dvcmQ="
                }
            }
        }
    }
}
```

Below is the example of bad request payload which should fail the policy
Here alice is trying to make `GET` request on path `api/v1/products` which is not allowed.
```json
{
    "attributes": {
        "request": {
            "http": {
                "method": "GET",
                "path": "/api/v1/products",
                "headers": {
                    "authorization": "Basic YWxpY2U6cGFzc3dvcmQ="
                }
            }
        }
    }
}
```

Below is the example of validation policy that restricts access to an endpoint based on a user’s role and permissions.
```yml
....

....
```